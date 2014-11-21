package mercury

import (
	"appengine"
	"appengine/datastore"
	"net/http"
	"time"
)

type kinder interface {
	kind() string
}

type empty struct {
}

func (e *empty) kind() string {
	return "empty"
}

func init() {
	http.HandleFunc("/fewKinds", fewKindsHandler)
}

func fewKindsHandler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	for i := 0; i < 1000; i++ {
		start := time.Now()
		e := &empty{}
		k := datastore.NewIncompleteKey(cxt, e.kind(), nil)
		datastore.Put(cxt, k, e)
		total := time.Now().Sub(start)
		cxt.Infof("Few Kinds: %v", total)
	}
}
