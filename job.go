package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorhill/cronexpr"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

func (j Job) String() string {
	return fmt.Sprintf("{name:'%s',command:'%s',cronExpr:'%s'}", j.Name, j.Command, j.CronExpr)
}

func UnpackJob(val []byte) (ret *Job, err error) {
	ret = &Job{}
	err = json.Unmarshal(val, ret)
	return
}

// job事件
type JobEvent struct {
	Type int
	Job  *Job
}

func (j JobEvent) String() string {
	return fmt.Sprintf("{type:%d,job:%s}", j.Type, j.Job)
}

func BuildEvent(eventType int, job *Job) *JobEvent {
	return &JobEvent{
		Type: eventType,
		Job:  job,
	}
}

// 任务执行计划
type JobSchedulePlan struct {
	Job        *Job
	Expr       *cronexpr.Expression // 解析好的表达式
	NextTime   time.Time            //下次执行时间
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (j JobSchedulePlan) String() string {
	return fmt.Sprintf("{job:%s,nextTime:%s}", j.Job, GetTimeString(j.NextTime))
}

type JobMgr struct {
	client *clientv3.Client
}

func (jm *JobMgr) Get(ctx context.Context, key string) (string, error) {
	fun := "JobMgr.Get -->"

	logrus.Infof("%s before get ...", fun)
	res, err := jm.client.Get(ctx, key, clientv3.WithPrevKV())
	logrus.Infof("%s get ...", fun)
	if err != nil {
		logrus.Errorf("%s client get key:%s err:%v", fun, key, err)
		return "", err
	}

	if res.Count <= 0 {

	}
	logrus.Infof("%s get %+v", fun, res)

	return "", nil
}

func (jm *JobMgr) Put(ctx context.Context, key, value string) (string, error) {
	fun := "JobMgr.Put -->"
	rs := ""

	res, err := jm.client.Put(ctx, key, value, clientv3.WithPrevKV())
	if err != nil {
		logrus.Errorf("%s client put key:%s err:%v", fun, key, err)
		return "", err
	}

	logrus.Infof("%s put res:%+v err:%v", fun, res, err)
	if res.PrevKv == nil {

		return "", nil
	}

	// TODO 老数据处理
	b, err := json.Marshal(res.PrevKv)
	if err != nil {
		logrus.Errorf("%s client put key:%s err:%v", fun, key, err)
		return "", err
	}

	rs = string(b)

	return rs, nil
}

func (jm *JobMgr) Delete(ctx context.Context, key string) (string, error) {

	return "", nil
}

func (jm *JobMgr) Watch(ctx context.Context, key string) chan *JobEvent {
	ch := make(chan *JobEvent, 1024)

	go jm.watch(ctx, key, ch)
	return ch
}

func (jm *JobMgr) watch(ctx context.Context, key string, ch chan *JobEvent) {
	fun := "JobMgr.watch -->"
	defer close(ch)

	logrus.Infof("%s start....", fun)
	for {
		select {
		case <-ctx.Done():
			// 退出
		default:
			rch := jm.client.Watch(
				clientv3.WithRequireLeader(ctx),
				key,
				clientv3.WithPrefix())

			logrus.Infof("%s key:%s T:%T len:%d resp:%+v time:%s", fun, key, rch, len(rch), rch, GetNowString())

			for wresp := range rch {
				logrus.Infof("%s range WatchChan resp:%+v", fun, wresp)
				if wresp.Created {
					logrus.Info("%s etcd watcher created", fun)
					continue
				}

				if wresp.Canceled {
					logrus.WithError(wresp.Err()).Error("watcher is canceled by etcd server")
					break
				}

				for _, evt := range wresp.Events {
					logrus.Infof("%s event:%+v", fun, evt)
					if evt.IsCreate() {

						logrus.WithFields(logrus.Fields{
							"key": string(evt.Kv.Key),
						}).Debug("key is created")
					} else if evt.IsModify() {

						logrus.WithFields(logrus.Fields{
							"key": string(evt.Kv.Key),
						}).Debug("key is modified")
					} else {
						// delete
						logrus.WithFields(logrus.Fields{
							"key": string(evt.Kv.Key),
						}).Debug("key is deleted")
					}

				}
			}

			logrus.Debugf("%s info", fun)
		}
	}

}

func (jm *JobMgr) KillJob(ctx context.Context, key string) error {
	//
	fun := "JobMgr.KillJob -->"
	lease, err := jm.client.Grant(ctx, 1)
	if err != nil {
		logrus.Errorf("%s Grant key:%s err:%v", fun, key, err)
		return err
	}

	_, err = jm.client.Put(ctx, key, "", clientv3.WithLease(lease.ID))
	if err != nil {
		logrus.Errorf("%s Put key:%s err:%v", fun, key, err)
	}

	return err
}

func (jm *JobMgr) initJobsByType(ctx context.Context, key string, typ int) (int64, error) {
	var (
		fun          = "JobMgr.initJobsByType -->"
		watchVersion int64
	)

	getResp, err := jm.client.Get(ctx, key, clientv3.WithPrefix())
	if err != nil {
		logrus.Errorf("%s Get key:%s err:%v", fun, key, err)
		return watchVersion, err
	}
	watchVersion = getResp.Header.Revision + 1

	// 初始化或者重启
	for _, v := range getResp.Kvs {
		logrus.Debugf("%s kv:%+v", fun, v)
		job := &Job{}
		err = json.Unmarshal(v.Value, job)
		if err != nil {
			logrus.Warnf("%s Unmarshal key:%s value:%s err:%v", fun, key, v.Value, err)
			continue
		}

		jobEvent := &JobEvent{
			Type: typ,
			Job:  job,
		}

		// job量大 阻塞
		go GScheduler.Push(jobEvent)
	}
	// 后面错误不报
	return watchVersion, nil
}

func (jm *JobMgr) WatchJobs(ctx context.Context, key string) error {
	fun := "JobMgr.WatchJobs -->"
	watchVersion, err := jm.initJobsByType(ctx, key, JOB_EVENT_SAVE)
	if err != nil {
		return err
	}

	logrus.Infof("%s watchVersion:%d", fun, watchVersion)
	go func(watchVersion int64) {
		watchChan := jm.client.Watch(ctx, key,
			clientv3.WithRev(watchVersion),
			clientv3.WithPrefix(),
		)

		for resp := range watchChan {
			for _, res := range resp.Events {
				logrus.Infof("%s Event res:%+v", fun, res)
				switch res.Type {
				case mvccpb.PUT:
					job := &Job{}
					err := json.Unmarshal(res.Kv.Value, job)
					if err != nil {
						logrus.Warnf("%s Unmarshal key:%s value:%s err:%v", fun, key, res.Kv.Value, err)
						continue
					}

					je := &JobEvent{
						Type: JOB_EVENT_SAVE,
						Job:  job,
					}
					go GScheduler.Push(je)
				case mvccpb.DELETE:
					logrus.Infof("%s DELETE key:%s value:%s", fun, key, res.Kv.Key)
				}
			}
		}
	}(watchVersion)

	return err
}

func (jm *JobMgr) WatchKillJobs(ctx context.Context, key string) error {
	fun := "JobMgr.WatchKillJobs -->"
	_, err := jm.initJobsByType(ctx, key, JOB_EVENT_DELETE)
	if err != nil {
		return err
	}

	go func() {
		watchChan := jm.client.Watch(ctx, key,
			clientv3.WithPrefix(),
		)

		for resp := range watchChan {
			for _, res := range resp.Events {
				logrus.Infof("%s Event res:%+v", fun, res)
				switch res.Type {
				case mvccpb.PUT:
				case mvccpb.DELETE:
				}
			}
		}
	}()

	return err
}

func (jm *JobMgr) NewJobLock(key string) *JobLock {
	return &JobLock{
		kv:      GClient.KV,
		lease:   GClient.Lease,
		jobName: key,
	}
}
