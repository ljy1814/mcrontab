package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTHandle(t *testing.T) {
	pl, err := ioutil.ReadFile("testdata/payload.json")
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

	if exp := "handle pull request id:191568743\n"; string(msg) != exp {
		t.Fatalf("expected message %q; got %q", exp, msg)
	}
}

func BenchmarkHandle9(b *testing.B) {
	b.StopTimer()

	fname := "testdata/payload.json"
	pl, err := ioutil.ReadFile(fname)
	if err != nil {
		b.Fatalf("could not read payload.json err:%v", err)
	}
	//logrus.Infof("TEST payload:%s \n%d", pl, b.N)

	uu := "http://127.0.0.1:8080"
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodPost, uu, bytes.NewReader(pl))
		if err != nil {
			b.Fatalf("could not create test request err:%v", err)
		}

		// logrus.Infof("req:%+v", req)
		rec := httptest.NewRecorder()
		b.StartTimer()
		handle(rec, req)
		//logrus.Infof("req:%d", i)
		b.StopTimer()
	}
}

func BenchmarkHandle0(b *testing.B) {
	b.StopTimer()

	fname := "testdata/payload.json"
	pl, err := ioutil.ReadFile(fname)
	if err != nil {
		b.Fatalf("could not read payload.json err:%v", err)
	}
	//logrus.Infof("TEST payload:%s \n%d", pl, b.N)
	uu := "http://127.0.0.1:8080/0"
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodPost, uu, bytes.NewReader(pl))
		if err != nil {
			b.Fatalf("could not create test request err:%v", err)
		}

		// logrus.Infof("req:%+v", req)
		rec := httptest.NewRecorder()
		b.StartTimer()
		handle(rec, req)
		//logrus.Infof("req:%d", i)
		b.StopTimer()
	}
}

func BenchmarkHandle1(b *testing.B) {
	b.StopTimer()

	fname := "testdata/payload.json"
	pl, err := ioutil.ReadFile(fname)
	if err != nil {
		b.Fatalf("could not read payload.json err:%v", err)
	}
	//logrus.Infof("TEST payload:%s \n%d", pl, b.N)
	uu := "http://127.0.0.1:8080/1"
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodPost, uu, bytes.NewReader(pl))
		if err != nil {
			b.Fatalf("could not create test request err:%v", err)
		}

		// logrus.Infof("req:%+v", req)
		rec := httptest.NewRecorder()
		b.StartTimer()
		handle(rec, req)
		//logrus.Infof("req:%d", i)
		b.StopTimer()
	}
}

func BenchmarkHandle2(b *testing.B) {
	b.StopTimer()

	fname := "testdata/payload.json"
	pl, err := ioutil.ReadFile(fname)
	if err != nil {
		b.Fatalf("could not read payload.json err:%v", err)
	}
	uu := "http://127.0.0.1:8080/2"
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodPost, uu, bytes.NewReader(pl))
		if err != nil {
			b.Fatalf("could not create test request err:%v", err)
		}

		// logrus.Infof("req:%+v", req)
		rec := httptest.NewRecorder()
		b.StartTimer()
		handle(rec, req)
		//logrus.Infof("req:%d", i)
		b.StopTimer()
	}
}

func BenchmarkHandle3(b *testing.B) {
	b.StopTimer()

	fname := "testdata/payload.json"
	pl, err := ioutil.ReadFile(fname)
	if err != nil {
		b.Fatalf("could not read payload.json err:%v", err)
	}

	uu := "http://127.0.0.1:8080/3"
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest(http.MethodPost, uu, bytes.NewReader(pl))
		if err != nil {
			b.Fatalf("could not create test request err:%v", err)
		}

		// logrus.Infof("req:%+v", req)
		rec := httptest.NewRecorder()
		b.StartTimer()
		handle(rec, req)
		//logrus.Infof("req:%d", i)
		b.StopTimer()
	}
}
