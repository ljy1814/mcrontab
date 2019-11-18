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
		ScheduleResultChan: make(chan *ExecResult, DefaultExecutorCount),
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

func stopServer() {
	GetHttpInstance().Stop()
}

func initLogger() {
	var (
		dir              = "./logs"
		buffersize int64 = 1024
		rotateSize int64 = 128 * 1024 * 1024
		maxLogFile int   = 5
	)
	GLogger = NewDemoFile(dir, buffersize, rotateSize, maxLogFile)

	//logrus.AddHook(logrusHook{})
	logrus.SetOutput(GLogger)
}
