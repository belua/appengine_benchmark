package aebench

import (
	"appengine"
	"appengine/datastore"
	"github.com/belua/httprouter"
	"net/http"
)

type empty struct {
}

type emptyBuilder struct{}

func (b *emptyBuilder) build(cxt appengine.Context) (*datastore.Key, interface{}) {
	return datastore.NewIncompleteKey(cxt, "empty", nil), &empty{}
}

func emptyHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &emptyBuilder{}
	putEntities(w, r, b)
}
