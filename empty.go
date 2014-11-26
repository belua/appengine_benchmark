package aebench

import (
	"github.com/belua/httprouter"
	"net/http"
)

type empty struct {
}

func (e *empty) kind() string {
	return "empty"
}

type emptyBuilder struct{}

func (b *emptyBuilder) build() kinder {
	return &empty{}
}

func emptyHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &emptyBuilder{}
	putKinderSequential(w, r, b)
}
