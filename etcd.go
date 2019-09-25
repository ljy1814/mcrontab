package main

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

func InitEtcd() (err error) {
	fun := "InitEtcd -->"
	config := clientv3.Config{
		Endpoints:   GEtcdCluster,
		DialTimeout: 3 * time.Second,
	}

	GClient, err = clientv3.New(config)
	if err != nil {
		logrus.Errorf("%s err:%v", fun, err)
		return
	}

	logrus.Infof("%s ok...", fun)

	return
}

type EtcdOp interface {
	Get(context.Context, string) (string, error)
	Put(context.Context, string, string) (string, error)
	Delete(context.Context, string) (string, error)
}
