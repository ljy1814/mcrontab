package main

import (
	"context"
	"runtime"

	"github.com/Sirupsen/logrus"
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initScheduler() error {
	return InitScheduler()
}

func initHttpServer() error {
	logrus.Info("init http server")

	port := "8080"
	hs := GetHttpInstance()
	hs.Init(port)

	go hs.Start()

	return nil
}

func initSignal() {
	signalProc()
}

func initEtcd() {
	InitEtcd()
}

func initJobMgr() {
	GJobMgr = &JobMgr{GClient}
}

func initExecutor() {
	GExecutor = &Executor{
		ScheduleResultChan: make(chan *ExecResult, 1024),
	}
}

func initWatcher() error {
	e1 := GJobMgr.WatchJobs(context.Background(), JOB_PREFIX_PUT)
	e2 := GJobMgr.WatchKillJobs(context.Background(), JOB_PREFIX_KILL)

	if e1 != nil {
		return e1
	}
	return e2
}
