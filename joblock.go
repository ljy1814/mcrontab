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
		logrus.Errorf("%s Lease KeepAlive key:%s err:%v", fun, key, err)
		return err
	}

	go func() {
		for {
			select {
			case r := <-leaseKeepAliveChan:
				if r == nil {
					return
				}
			}
		}
	}()

	txn := j.kv.Txn(ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(key))

	txnResp, err := txn.Commit()
	if err != nil {
		logrus.Errorf("%s Lease txn commit key:%s err:%v", fun, key, err)
		return err
	}

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
