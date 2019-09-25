package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
)

func signalProc() {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGALRM, syscall.SIGTERM, syscall.SIGUSR1)

	sig := <-c

	logrus.Warnf("Signal received: %v", sig)

	//httpserver.HttpListener.Close()

	time.Sleep(100 * time.Millisecond)

}
