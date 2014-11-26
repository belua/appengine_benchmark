package aebench

import (
	"net/http"
	"github.com/belua/httprouter"
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

func emptyDelHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	delKind(w, r, (&empty{}).kind())
}
