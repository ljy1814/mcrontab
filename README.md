### crontab项目

```
ENDPOINTS='http://10.0.2.15:4001,http://10.0.2.15:5001,http://10.0.2.15:6001'
./etcdctl --endpoints="$ENDPOINTS" put  /jobs/put '{"name":"test" ,"command" : "echo hello" , "cronExpr" : "*/1 * * * * " }'

```

```

curl -X POST -L http://127.0.0.1:4001/v3beta/kv/put  -d '{"key": "QUFB", "value": "QkJCQg=="}' | python -mjson.tool

{
    "header": {
        "cluster_id": "7424915042368147604",
        "member_id": "8687557477222935394",
        "raft_term": "3733",
        "revision": "612"
    }
}

curl -X POST -L http://127.0.0.1:4001/v3beta/kv/range  -d '{"key" : "QUFB"}' | python -mjson.tool

{
    "count": "1",
    "header": {
        "cluster_id": "7424915042368147604",
        "member_id": "8687557477222935394",
        "raft_term": "3733",
        "revision": "618"
    },
    "kvs": [
        {
            "create_revision": "611",
            "key": "QUFB",
            "mod_revision": "612",
            "value": "QkJCQg==",
            "version": "2"
        }
    ]
}

```

raft
```

https://ggaaooppeenngg.github.io/
https://etcd.io/docs/v3.4.0/demo/#lease
http://lday.me/2017/02/01/0003_seri-stm-etcd3/
https://github.com/dollarkillerx/GO-Distributed-Task-Scheduling

```

python base64
```python
import base64
base64.b64encode('AA')
```

etcdctl watch  /test/ok

#### 如果监听子节点
etcdctl watch  /test/ok --prefix

etcdctl lease grant 40
lease 4e5e5b853f528859 granted with TTL(40s)

etcdctl put --lease=4e5e5b853f528859 /test/ok/first xx
OK

etcdctl lease revoke 4e5e5b853f5286cc
lease 4e5e5b853f5286cc revoked

etcdctl lease keep-alive 4e5e5b853f52892b
lease 4e5e5b853f52892b keepalived with TTL(40)
lease 4e5e5b853f52892b keepalived with TTL(40)



### 其它资源
```
https://github.com/soimort/you-get
```

```

添加time_wait   close_wait连接数
netstat -n | awk '/^tcp/ {++S[$NF]} END {for(a in S) print a, S[a]}'

```

```
https://github.com/chai2010/advanced-go-programming-book/blob/master/ch1-basic/ch1-03-array-string-and-slice.md
https://zboya.github.io/post/go_scheduler/
https://www.shkuro.com/books/2019-mastering-distributed-tracing/
```

#### 虚拟内存
```
http://blog.coderhuo.tech/2017/10/12/Virtual_Memory_C_strings_proc/
```

任务调度开源项目
```
https://github.com/joyent/containerpilot.git
```

logrus
```
AddHook 对特定的日志加某种操作,不改变日志输出
SetOutput 设置输出目的地
```
