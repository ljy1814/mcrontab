package main

import (
	"context"
	"fmt"

	"github.com/Sirupsen/logrus"
)

func main() {
	initEnv()
	initScheduler()
	initHttpServer()
	initEtcd()
	initJobMgr()
	initExecutor()
	initWatcher()

	logrus.Infof("init ok...")
	GJobMgr.Put(context.Background(), "/jobs", "jobtest")
	logrus.Info("------")
	GJobMgr.Get(context.Background(), "/jobs")

	GJobMgr.Watch(context.Background(), "/jobs")
	for i := 0; i < 3; i++ {
		GJobMgr.Put(context.Background(), "/jobs", fmt.Sprintf("jobtest:%d", i))
	}
	initSignal()
	logrus.Info("server quit")
}
