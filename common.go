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

func init() {
	http.HandleFunc("/emptySequential", emptySequentialHandler)
	http.HandleFunc("/emptyParallel", emptyParallelHandler)
	http.HandleFunc("/delEmpty", delEmptyHandler)
}

func putKinder(cxt appengine.Context, entity kinder, count int) error {
	start := time.Now()
	key := datastore.NewIncompleteKey(cxt, entity.kind(), nil)
	if _, err := datastore.Put(cxt, key, entity); err != nil {
		return err
	}
	total := time.Now().Sub(start)
	cxt.Infof("%d %s Single Put: %v", count, entity.kind(), total)
	return nil
}
