#### sync.pool test

简单测试
```shell
curl -X POST -d @testdata/payload.json http://127.0.0.1:8080/
```

压测
```shell
go test -bench=. 
go test -bench=. -benchmem
```
