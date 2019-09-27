package main

import (
	"net"
	"os"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

var (
	onceHttp     sync.Once
	httpInstance *HttpServer
)

type HttpServer struct {
	//Route  *httprouter.Router
	Port   string
	Server *http.Server
}

func GetHttpInstance() *HttpServer {
	logrus.Info("get http instance")
	onceHttp.Do(
		func() {
			if httpInstance == nil {
				httpInstance = &HttpServer{}
			}
		})

	return httpInstance
}

func (hs *HttpServer) Init(port string) error {
	var (
		fun = "HttpServer.Init -->"
	)
	logrus.Infof("%s start", fun)

	router := httprouter.New()
	router.GET("/", Index)
	hs.Port = port

	//	hs.ServerMux = http.ewServeMux()
	hs.Server = &http.Server{
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return nil
}

func (hs *HttpServer) Start() error {
	fun := "HttpServer.Start -->"
	curPid := os.Getpid()
	osProcess := os.Process{Pid: curPid}
	logrus.Infof("%s pid:%d start...", fun, curPid)

	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("%s Panic pid:%d err:%v stack:%s", fun, curPid, err, debug.Stack())
		}
	}()

	hl, err := net.Listen("tcp", ":"+hs.Port)
	if err != nil {
		logrus.Errorf("%s Listen pid:%d err:%v", fun, curPid, err)
		osProcess.Signal(syscall.SIGINT)
		return err
	}
	defer hl.Close()

	err = hs.Server.Serve(hl)
	if err != nil {
		logrus.Errorf("%s Serve pid:%d err:%v", fun, curPid, err)
		return err
	}
	return err
}

func Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	logrus.Infof("Index .....req:%+v", r)
	w.Write([]byte("OK\r\n"))
	return
}
