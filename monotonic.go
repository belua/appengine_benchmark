package aebench

import (
	"appengine"
	"appengine/datastore"
	"github.com/belua/httprouter"
	"net/http"
)

// Entity with single indexable field which increases monotonically

type MonoIndex struct {
	Index int64
}

type MonoIndexBuilder struct {
	indexVal int64
}

func (b *MonoIndexBuilder) build(cxt appengine.Context) (*datastore.Key, interface{}) {
	index := b.indexVal
	b.indexVal++
	return datastore.NewIncompleteKey(cxt, "MonoIndex", nil), MonoIndex{Index: index}
}

func monoIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &MonoIndexBuilder{}
	putEntities(w, r, b)
}
