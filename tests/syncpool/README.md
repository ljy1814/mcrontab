#### sync.pool test

简单测试
```shell
curl -X POST -d @testdata/payload.json http://127.0.0.1:8080/
```

压测

handle0
```shell
go test -bench=. 
goos: linux
goarch: amd64
BenchmarkHandle 	    3000	    454256 ns/op
PASS
ok  	_/home/vagrant/syncpool	1.452s

go test -bench=. -benchmem
goos: linux
goarch: amd64
BenchmarkHandle 	    3000	    459204 ns/op	   84423 B/op	     541 allocs/op
PASS
ok  	_/home/vagrant/syncpool	1.468s
```

handle1
```shell
go test -bench=. 
goos: linux
goarch: amd64
BenchmarkHandle 	    3000	    460155 ns/op
PASS
ok  	_/home/vagrant/syncpool	1.472s

go test -bench=. -benchmem
goos: linux
goarch: amd64
BenchmarkHandle 	    3000	    455378 ns/op	   84423 B/op	     541 allocs/op
PASS
ok  	_/home/vagrant/syncpool	1.460s
```

handle2
```shell
go test -bench=. 
goos: linux
goarch: amd64
BenchmarkHandle 	    2000	    532779 ns/op
PASS
ok  	_/home/vagrant/syncpool	1.168s

go test -bench=. -benchmem
goos: linux
goarch: amd64
BenchmarkHandle 	    2000	    525204 ns/op	   84424 B/op	     541 allocs/op
PASS
ok  	_/home/vagrant/syncpool	1.154s
```

handle3
```shell
go test -bench=. 
goos: linux
goarch: amd64
BenchmarkHandle 	    3000	    539890 ns/op
PASS
ok  	_/home/vagrant/syncpool	1.725s

go test -bench=. -benchmem
goos: linux
goarch: amd64
BenchmarkHandle 	    3000	    513894 ns/op	   94197 B/op	     883 allocs/op
PASS
ok  	_/home/vagrant/syncpool	1.643s
```

handle4
```shell
go test -bench=. 
goos: linux
goarch: amd64
BenchmarkHandle 	   10000	    224927 ns/op
PASS
ok  	_/home/vagrant/syncpool	2.387s

go test -bench=. -benchmem
goos: linux
goarch: amd64
BenchmarkHandle 	    5000	    228730 ns/op	   68557 B/op	     246 allocs/op
PASS
ok  	_/home/vagrant/syncpool	1.231s
```

```
go test -bench=. -benchmem -memprofile=mem.pb.gz
```
