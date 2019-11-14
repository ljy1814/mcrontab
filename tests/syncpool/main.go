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
	fmt.Fprintf(w, "handle0 pull request id:%d\n", *data.PullRequest.ID)
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
	fmt.Fprintf(w, "handle1 pull request id:%d\n", *data.PullRequest.ID)
}

func handle2(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprintf(w, "handle2 pull request id:%d\n", *data.PullRequest.ID)
}

func handle3(w http.ResponseWriter, r *http.Request) {
	data := prPool.Get().(*github.PullRequestEvent)
	defer prPool.Put(data)

	*data = github.PullRequestEvent{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logrus.Errorf("could not decode request: %v", err)
		http.Error(w, "could not decode request", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "handle3 pull request id:%d\n", *data.PullRequest.ID)
}

type (
	pullRequest struct {
		PullRequest struct {
			ID *int `json:"id,omitempty"`
		} `json:"pull_request,omitempty"`
	}
)

var (
	plPool = sync.Pool{
		New: func() interface{} { return new(pullRequest) },
	}
)

func handle(w http.ResponseWriter, r *http.Request) {
	data := plPool.Get().(*pullRequest)
	defer plPool.Put(data)

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logrus.Errorf("could not decode request: %v", err)
		http.Error(w, "could not decode request", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "handle pull request id:%d\n", *data.PullRequest.ID)
}

func main() {
	http.HandleFunc("/", handle3)
	http.HandleFunc("/0", handle0)
	http.HandleFunc("/1", handle1)
	http.HandleFunc("/2", handle2)
	http.HandleFunc("/3", handle3)
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}
