package main

import (
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
