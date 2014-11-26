package aebench

import (
	"github.com/belua/httprouter"
	"net/http"
)

// Entity with single indexable field which increases monotonically

type MonoIndex struct {
	Index int64
}

func (i *MonoIndex) kind() string {
	return "MonoIndex"
}

type MonoIndexBuilder struct {
	indexVal int64
}

func (b *MonoIndexBuilder) build() kinder {
	index := b.indexVal
	b.indexVal++
	return &MonoIndex{Index: index}
}

func monoIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &MonoIndexBuilder{}
	putKinderSequential(w, r, b)
}
