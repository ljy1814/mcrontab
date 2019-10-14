package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorhill/cronexpr"
)

type Scheduler struct {
	//
	jobEventChan chan *JobEvent
	// job执行计划表
	jobEventPlanMap map[string]*SchedulePlan
	// 正在执行的job表
	jobExecMap map[string]*SchedulePlan
	lock       *sync.Mutex
}

func (s *Scheduler) scheduleLoop() {
	var (
		timeAfter      = s.trySchedule()
		schedulerTimer = time.NewTimer(timeAfter)
	)
	for {
		select {
		case v := <-s.jobEventChan:
			logrus.Info(v)
			// 任务更新
			s.handle(context.Background(), v)
		case <-schedulerTimer.C:
			// 休眠结束
		case result := <-GExecutor.ScheduleResultChan:
			// 任务执行完毕
			//logrus.Info(result)
			s.processScheduleResult(result)
		}

		timeAfter = s.trySchedule()
		schedulerTimer.Reset(timeAfter)
	}
}

func (s *Scheduler) trySchedule() (timeAfter time.Duration) {
	if len(s.jobEventPlanMap) == 0 {
		return time.Second
	}

	now := time.Now()
	var nearestTime *time.Time
	for _, v := range s.jobEventPlanMap {
		if v.NextTime.Before(now) || v.NextTime.Equal(now) {
			s.tryStartJob(v)
			v.NextTime = v.Expr.Next(now)
		}

		if nearestTime == nil || v.NextTime.Before(*nearestTime) {
			nearestTime = &v.NextTime
		}
	}

	return nearestTime.Sub(now)
}

func (s *Scheduler) Push(je *JobEvent) {
	s.jobEventChan <- je
}

func (s *Scheduler) tryStartJob(plan *SchedulePlan) {
	fun := "Scheduler.tryStartJob -->"

	// 任务还在执行则跳过本次调度
	if _, ok := s.jobExecMap[plan.Job.Name]; ok {
		logrus.Warnf("%s job:%s executeing", fun, plan.Job)
		return
	}

	// 加入执行表
	s.jobExecMap[plan.Job.Name] = plan
	GExecutor.ExecJob(plan)
}

func (s *Scheduler) handle(ctx context.Context, je *JobEvent) {
	var (
		fun = "Scheduler.handle -->"
	)
	s.lock.Lock()
	defer s.lock.Unlock()

	switch je.Type {
	case JOB_EVENT_SAVE:
		logrus.Infof("%s SAVE job:%s", fun, je)
		jobSchedulePlan, err := buildJobSchedulePlan(ctx, je.Job)
		if err != nil {
			logrus.Errorf("%s buildJobSchedulePlan ctx:%v job:%s err:%v", fun, ctx, je, err)
			return
		}
		s.jobEventPlanMap[je.Job.Name] = jobSchedulePlan

	case JOB_EVENT_DELETE:
		logrus.Infof("%s DELETE job:%s", fun, je)
		if _, ok := s.jobEventPlanMap[je.Job.Name]; ok {
			delete(s.jobEventPlanMap, je.Job.Name)
		}

	case JOB_EVENT_KILL:
		jobSchedulePlan, ok := s.jobExecMap[je.Job.Name]
		if !ok {
			return
		}

		logrus.Infof("%s KILL job:%s", fun, je)
		jobSchedulePlan.cancelFunc()
		// 重置强杀的任务
		jobSchedulePlan.ctx, jobSchedulePlan.cancelFunc = context.WithCancel(ctx)
	}
}

func (s *Scheduler) processScheduleResult(result *ExecResult) {
	fun := "Scheduler.processScheduleResult -->"

	s.lock.Lock()
	delete(s.jobExecMap, result.JobPlan.Job.Name)
	s.lock.Unlock()

	time.Sleep(20 * time.Millisecond)
	logrus.Infof("%s result:%s", fun, result)

	// TODO 记录日志
}

func buildJobSchedulePlan(ctx context.Context, job *Job) (*SchedulePlan, error) {
	fun := "buildJobSchedulePlan -->"
	expr, err := cronexpr.Parse(job.CronExpr)
	if err != nil {
		logrus.Warnf("%s ctx:%v cronexpr Parse job:%s err:%v", fun, ctx, job, err)
		return nil, err
	}

	ctx1, cancelFunc := context.WithCancel(ctx)

	return &SchedulePlan{
		Job:        job,
		Expr:       expr,
		NextTime:   expr.Next(time.Now()),
		ctx:        ctx1,
		cancelFunc: cancelFunc,
	}, err
}

var (
	GScheduler *Scheduler
)

func InitScheduler() (err error) {
	GScheduler = &Scheduler{
		jobEventChan: make(chan *JobEvent, DefaultSchedulerCount),
		// job执行计划表
		jobEventPlanMap: make(map[string]*SchedulePlan),
		// 正在执行的job表
		jobExecMap: make(map[string]*SchedulePlan),
		lock:       &sync.Mutex{},
	}

	go GScheduler.scheduleLoop()

	return
}

type SchedulePlan struct {
	Job        *Job
	Expr       *cronexpr.Expression
	NextTime   time.Time
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (sp SchedulePlan) String() string {
	return fmt.Sprintf("{job:'%s',nextTime:'%s'}", sp.Job, GetTimeString(sp.NextTime))
}
