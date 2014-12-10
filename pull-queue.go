package aebench

import (
	"appengine/taskqueue"
	"github.com/belua/httprouter"
	"net/http"
)

const pullQueueName = "perf-pull-queue"

func pullQueueHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := &taskqueue.Task{
		Payload: []byte("hello pull-queue"),
		Method:  "PULL",
	}
	addToQueue(w, r, pullQueueName, t)
}
