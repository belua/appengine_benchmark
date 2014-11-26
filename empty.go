package mercury

import (
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

func emptyHandler(w http.ResponseWriter, r *http.Request) {
	b := &emptyBuilder{}
	putKinderSequential(w, r, b)
}

func emptyDelHandler(w http.ResponseWriter, r *http.Request) {
	delKind(w, r, (&empty{}).kind())
}
