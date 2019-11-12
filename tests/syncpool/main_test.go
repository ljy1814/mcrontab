package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkHandle(b *testing.B) {
	b.StopTimer()

	fname := "testdata/payload.json"
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

		// logrus.Infof("req:%+v", req)
		rec := httptest.NewRecorder()
		handle(rec, req)
		//logrus.Infof("req:%+v", req)
		b.StopTimer()
	}
}
