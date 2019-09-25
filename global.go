package main

import "go.etcd.io/etcd/clientv3"

var (
	GQuitChan = make(chan int)
)
var (
	GExecutor *Executor = &Executor{}
)

var (
	GClient *clientv3.Client

	GJobMgr *JobMgr
)

var (
	//GEtcdCluster []string = []string{"http://127.0.0.1:2041", "http://127.0.0.1:2051", "http://127.0.0.1:2061"}
	GEtcdCluster []string = []string{"http://127.0.0.1:4001", "http://127.0.0.1:5001", "http://127.0.0.1:6001"}
)
