package aebench

import (
	"appengine"
	"appengine/datastore"
	"github.com/belua/httprouter"
	"net/http"
)

// Entity with single indexable field which increases monotonically

type StringList struct {
	Strings []string
}

type StringListBuilder struct {
}

func (b *StringListBuilder) build(cxt appengine.Context) (*datastore.Key, interface{}) {
	strings := make([]string, 0)
	for i := 0; i < 1000; i++ {
		strings = append(strings, "Moderately Long String, I think is is probably roughly enough")
	}
	return datastore.NewIncompleteKey(cxt, "stringList", nil), &StringList{strings}
}

func stringListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &StringListBuilder{}
	putEntities(w, r, b)
}
