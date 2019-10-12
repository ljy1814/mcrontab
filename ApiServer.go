package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
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
	hs.Port = port
	router := hs.initHandler()

	//	hs.ServerMux = http.ewServeMux()
	hs.Server = &http.Server{
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return nil
}

func (hs *HttpServer) initHandler() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/close", Close)
	router.POST("/job/create", Create)

	return router
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

func (hs *HttpServer) Stop() error {
	logrus.Infof("HttpServer.Stop pid:%d ...", os.Getpid())
	return hs.Server.Shutdown(context.Background())
}

func Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	logrus.Infof("Index .....req:%+v", r)
	w.Write([]byte("OK\r\n"))
	return
}

func Close(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	logrus.Infof("Close .....req:%+v", r)
	time.Sleep(2 * time.Second)
	r.Close = true
	return
}

func Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fun := "HttpServer.Create -->"
	logrus.Infof("%s .....req:%+v", fun, r)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("%s req:%+v err:%v", fun, r, err)
		return
	}

	j := &Job{}
	err = json.Unmarshal(reqBody, j)
	if err != nil {
		logrus.Errorf("%s reqBody:%s err:%v", fun, reqBody, err)
		return
	}

	logrus.Infof("%s reqBody:%s job:%s err:%v", fun, reqBody, j, err)

	ctx := context.Background()
	ret, err := GJobMgr.Put(ctx, GetJobCreateKey(j.Name), string(reqBody))
	if err != nil {
		logrus.Errorf("%s PUT reqBody:%s err:%v", fun, reqBody, err)
		return
	}

	logrus.Infof("%s reqBody:%s resp:%s err:%v", fun, reqBody, ret, err)
	return
}

func Delete(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fun := "HttpServer.Create -->"
	logrus.Infof("%s .....req:%+v", fun, r)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("%s req:%+v err:%v", fun, r, err)
		return
	}

	type DeleteJob struct {
		Name string
	}
	dj := &DeleteJob{}
	err = json.Unmarshal(reqBody, dj)
	if err != nil {
		logrus.Errorf("%s reqBody:%s err:%v", fun, reqBody, err)
		return
	}

	logrus.Infof("%s reqBody:%s job:%s err:%v", fun, reqBody, dj, err)

	ctx := context.Background()
	ret, err := GJobMgr.Put(ctx, GetJobCreateKey(dj.Name), string(reqBody))
	if err != nil {
		logrus.Errorf("%s PUT reqBody:%s err:%v", fun, reqBody, err)
		return
	}

	logrus.Infof("%s reqBody:%s resp:%s err:%v", fun, reqBody, ret, err)
	return
}
