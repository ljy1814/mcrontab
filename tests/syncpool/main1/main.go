package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/src-d/go-github/github"
)

var (
	prPool = sync.Pool{
		New: func() interface{} { return new(github.PullRequestEvent) },
	}
)

func handle(w http.ResponseWriter, r *http.Request) {
	data := prPool.Get().(*github.PullRequestEvent)
	defer prPool.Put(data)

	*data = github.PullRequestEvent{}

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
