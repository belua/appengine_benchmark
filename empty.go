package mercury

import (
	"appengine"
	"appengine/datastore"
	"net/http"
	"time"
)

type empty struct {
}

func (e *empty) kind() string {
	return "empty"
}

const emptyCount = 1000

func emptySequentialHandler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	outerStart := time.Now()
	for i := 0; i < emptyCount; i++ {
		if err := putKinder(cxt, &empty{}, i); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			cxt.Infof("%s", err.Error())
			return
		}
	}
	outerTotal := time.Now().Sub(outerStart)
	cxt.Infof("Few Kinds %d Puts: %v", emptyCount, outerTotal)
}

func emptyParallelHandler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	outerStart := time.Now()
	complete := make(chan bool, emptyCount)
	for i := 0; i < emptyCount; i++ {
		go func(count int) {
			if err := putKinder(cxt, &empty{}, count); err != nil {
				cxt.Infof("%s", err.Error())
			}
			complete <- true
		}(i)
	}
	for i := 0; i < emptyCount; i++ {
		<-complete
	}
	outerTotal := time.Now().Sub(outerStart)
	cxt.Infof("Few Kinds %d Puts: %v", emptyCount, outerTotal)
}

func delEmptyHandler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	q := datastore.NewQuery((&empty{}).kind()).KeysOnly().Limit(500)
	var keys []*datastore.Key
	var err error
	for keys == nil || len(keys) > 0 {
		keys, err = q.GetAll(cxt, make([]*empty, 0))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			cxt.Infof("%s", err.Error())
			return
		}
		cxt.Infof("Deleting Empties: %d", len(keys))
		if err := datastore.DeleteMulti(cxt, keys); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			cxt.Infof("%s", err.Error())
			return
		}
		cxt.Infof("Sucessfully Deleted Empties: %d", len(keys))
	}
}
