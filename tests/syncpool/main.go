package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/src-d/go-github/github"
)

func handle0(w http.ResponseWriter, r *http.Request) {
	var data github.PullRequestEvent

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logrus.Errorf("could not decode request: %v", err)
		http.Error(w, "could not decode request", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "pull request id:%d\n", *data.PullRequest.ID)
}

var (
	prPool = sync.Pool{
		New: func() interface{} { return new(github.PullRequestEvent) },
	}
)

func handle1(w http.ResponseWriter, r *http.Request) {
	data := prPool.Get().(*github.PullRequestEvent)
	defer prPool.Put(data)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logrus.Errorf("could not decode request: %v", err)
		http.Error(w, "could not decode request", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "pull request id:%d\n", *data.PullRequest.ID)
}

func handle(w http.ResponseWriter, r *http.Request) {
	data := prPool.Get().(*github.PullRequestEvent)
	defer prPool.Put(data)

	if data.PullRequest != nil {
		data.PullRequest.ID = nil
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logrus.Errorf("could not decode request: %v", err)
		http.Error(w, "could not decode request", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "pull request id:%d\n", *data.PullRequest.ID)
}

func main() {
	http.HandleFunc("/", handle)
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}
