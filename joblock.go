package main

import (
	"context"
	"errors"

	"github.com/Sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

type JobLock struct {
	kv         clientv3.KV
	lease      clientv3.Lease
	leaseId    clientv3.LeaseID
	cancelFunc context.CancelFunc
	jobName    string
	isLocked   bool
}

func (j *JobLock) TryLock(ctx0 context.Context, key string) error {
	fun := "JobLock.TryLock -->"
	ctx, ctlFunc := context.WithCancel(ctx0)

	leaseResp, err := j.lease.Grant(ctx, 5)
	if err != nil {
		logrus.Errorf("%s Lease Grant key:%s err:%v", fun, key, err)
		return err
	}

	leaseId := leaseResp.ID
	defer func() {
		// 取消租约
		if err != nil {
			ctlFunc()
			j.lease.Revoke(ctx0, leaseId)
		}
	}()

	leaseKeepAliveChan, err := j.lease.KeepAlive(ctx, leaseId)
	if err != nil {
		logrus.Errorf("%s Lease KeepAlive leaseId:%d key:%s err:%v", fun, leaseId, key, err)
		return err
	}

	go func() {
		logrus.Infof("%s go func leaseId:%d starting...", fun, leaseId)
		defer logrus.Infof("%s go func leaseId:%d ending...", fun, leaseId)
		for {
			select {
			case r := <-leaseKeepAliveChan:
				// 最后收到空resp
				logrus.Infof("%s ||| leaseId:%d resp:%+v", fun, leaseId, r)
				if r == nil {
					return
				}
			case <-ctx0.Done():
				logrus.Infof("%s ctx0 go func leaseId:%d done...", fun, leaseId)
				// 这个分支收到消
			case <-ctx.Done():
				logrus.Infof("%s ctx go func leaseId:%d done...", fun, leaseId)
			}
		}
	}()

	txn := j.kv.Txn(ctx)
	t1 := clientv3.CreateRevision(key)
	logrus.Infof("%s ------leaseId:%d key:%s CreateRevision:%+v", fun, leaseId, key, t1)
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(key))

	txnResp, err := txn.Commit()
	if err != nil {
		logrus.Errorf("%s Lease:%d txn commit key:%s err:%v", fun, leaseId, key, err)
		return err
	}
	logrus.Infof("%s Commit leaseId:%d key:%s txnResp:%+v err:%v", fun, leaseId, key, txnResp, err)
	if !txnResp.Succeeded {
		err = errors.New("Lock already get by others")
		return err
	}

	j.leaseId = leaseId
	j.cancelFunc = ctlFunc
	j.isLocked = true

	return err
}

func (j *JobLock) UnLock(ctx context.Context) {
	if j.isLocked {
		j.cancelFunc()
		j.lease.Revoke(ctx, j.leaseId)
	}
}
