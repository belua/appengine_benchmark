package aebench

import (
	"appengine"
	"appengine/datastore"
	"github.com/belua/httprouter"
	"net/http"
)

const groupKeyVal = 5577006791947779410

type group struct {
}

type groupBuilder struct{}

func (b *groupBuilder) build(cxt appengine.Context) (*datastore.Key, interface{}) {
	groupKey := datastore.NewKey(cxt, "GroupKey", "", groupKeyVal, nil)
	return datastore.NewIncompleteKey(cxt, "Group", groupKey), &group{}
}

func groupHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &groupBuilder{}
	putEntitySequential(w, r, b)
}
