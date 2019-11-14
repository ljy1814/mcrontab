package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTHandle(t *testing.T) {
	pl, err := ioutil.ReadFile("../testdata/payload.json")
	if err != nil {
		t.Fatalf("could not read payload.json: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080", bytes.NewReader(pl))
	if err != nil {
		t.Fatalf("could not create test request: %v", err)
	}
	rec := httptest.NewRecorder()
	handle(rec, req)
	res := rec.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code %s", res.Status)
	}
	defer res.Body.Close()

	msg, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read result payload: %v", err)
	}

	if exp := "pull request id:191568743\n"; string(msg) != exp {
		t.Fatalf("expected message %q; got %q", exp, msg)
	}
}

/*
09:28 $ go test -bench=Handle -benchmem
goos: linux
goarch: amd64
pkg: github.com/ljy1814/crontab/master/tests/syncpool/main1
BenchmarkHandle 	    3000	    510304 ns/op	   94197 B/op	     883 allocs/op
PASS
ok  	github.com/ljy1814/crontab/master/tests/syncpool/main1	1.633s
*/
func BenchmarkHandle(b *testing.B) {
	b.StopTimer()

	fname := "../testdata/payload.json"
	pl, err := ioutil.ReadFile(fname)
	if err != nil {
		b.Fatalf("could not read payload.json err:%v", err)
	}

	uu := "http://127.0.0.1:8080"
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodPost, uu, bytes.NewReader(pl))
		if err != nil {
			b.Fatalf("could not create test request err:%v", err)
		}

		rec := httptest.NewRecorder()
		b.StartTimer()
		handle(rec, req)
		b.StopTimer()
	}
}
