#!/bin/bash
# Filename : etcd.sh
# Description : 
# Author : yajin
# Email : yajin160305@gmail.com
# Version
# Change : 
# Date : 2019-09-22 14:23:24

nohup ./etcd --name 'etcdNode1' \
    --listen-client-urls 'http://127.0.0.1:4001,http://10.0.2.15:4001' \
    --advertise-client-urls 'http://10.0.2.15:4001' \
    --listen-peer-urls 'http://10.0.2.15:4000' \
    --initial-advertise-peer-urls 'http://10.0.2.15:4000' \
    --initial-cluster-token "etcd-cluster" \
    --initial-cluster "etcdNode1=http://10.0.2.15:4000,etcdNode2=http://10.0.2.15:5000,etcdNode3=http://10.0.2.15:6000" \
    --data-dir '/home/arch/logs/etcd/node1/etcdNode1' 1>>/home/arch/logs/etcd/node1.log 2>&1 &


nohup ./etcd --name 'etcdNode2' \
    --listen-client-urls 'http://127.0.0.1:5001,http://10.0.2.15:5001' \
    --advertise-client-urls 'http://10.0.2.15:5001' \
    --listen-peer-urls 'http://10.0.2.15:5000' \
    --initial-advertise-peer-urls 'http://10.0.2.15:5000' \
    --initial-cluster-token "etcd-cluster" \
    --initial-cluster "etcdNode1=http://10.0.2.15:4000,etcdNode2=http://10.0.2.15:5000,etcdNode3=http://10.0.2.15:6000" \
    --data-dir '/home/arch/logs/etcd/node2/etcdNode2' 1>>/home/arch/logs/etcd/node2.log 2>&1 &

nohup ./etcd --name 'etcdNode3' \
    --listen-client-urls 'http://127.0.0.1:6001,http://10.0.2.15:6001' \
    --advertise-client-urls 'http://10.0.2.15:6001' \
    --listen-peer-urls 'http://10.0.2.15:6000' \
    --initial-advertise-peer-urls 'http://10.0.2.15:6000' \
    --initial-cluster-token "etcd-cluster" \
    --initial-cluster "etcdNode1=http://10.0.2.15:4000,etcdNode2=http://10.0.2.15:5000,etcdNode3=http://10.0.2.15:6000" \
    --data-dir '/home/arch/logs/etcd/node3/etcdNode3' 1>>/home/arch/logs/etcd/node3.log 2>&1 &
