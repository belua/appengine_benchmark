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
	http.HandleFunc("/fewKindsSequential", fewKindsSequentialHandler)
	http.HandleFunc("/fewKindsParallel", fewKindsParallelHandler)
	http.HandleFunc("/delEmpty", delEmptyHandler)
}

const fewKindsCount = 1000
const delLimit = 500

func fewKindsSequentialHandler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	outerStart := time.Now()
	for i := 0; i < fewKindsCount; i++ {
		if err := putEmpty(cxt, i); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			cxt.Infof("%s", err.Error())
			return
		}
	}
	outerTotal := time.Now().Sub(outerStart)
	cxt.Infof("Few Kinds %d Puts: %v", fewKindsCount, outerTotal)
}

func fewKindsParallelHandler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	outerStart := time.Now()
	complete := make(chan bool, fewKindsCount)
	for i := 0; i < fewKindsCount; i++ {
		go func(count int) {
			if err := putEmpty(cxt, count); err != nil {
				cxt.Infof("%s", err.Error())
			}
			complete <- true
		}(i)
	}
	for i := 0; i < fewKindsCount; i++ {
		<-complete
	}
	outerTotal := time.Now().Sub(outerStart)
	cxt.Infof("Few Kinds %d Puts: %v", fewKindsCount, outerTotal)
}

func putEmpty(cxt appengine.Context, count int) error {
	start := time.Now()
	e := &empty{}
	k := datastore.NewIncompleteKey(cxt, e.kind(), nil)
	if _, err := datastore.Put(cxt, k, e); err != nil {
		return err
	}
	total := time.Now().Sub(start)
	cxt.Infof("%d Few Kinds Single Put: %v", count, total)
	return nil
}

func delEmptyHandler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	q := datastore.NewQuery((&empty{}).kind()).KeysOnly().Limit(delLimit)
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
