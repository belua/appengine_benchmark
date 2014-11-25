package mercury

import (
	"appengine"
	"appengine/datastore"
	"net/http"
	"time"
)

const operationCount = 200

type kinder interface {
	kind() string
}

type kinderBuilder interface {
	build() kinder
}

func init() {
	http.HandleFunc("/emptyParallel", emptyParallelHandler)
	http.HandleFunc("/empty", emptyHandler)
	http.HandleFunc("/emptyDel", emptyDelHandler)
	http.HandleFunc("/oneIndex", oneIndexHandler)
	http.HandleFunc("/oneIndexDel", oneIndexDelHandler)
	http.HandleFunc("/twoIndex", oneIndexHandler)
	http.HandleFunc("/twoIndexDel", oneIndexDelHandler)
	http.HandleFunc("/threeIndex", oneIndexHandler)
	http.HandleFunc("/threeIndexDel", oneIndexDelHandler)
	http.HandleFunc("/fourIndex", oneIndexHandler)
	http.HandleFunc("/fourIndexDel", oneIndexDelHandler)
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

func delKey(cxt appengine.Context, key *datastore.Key, count int) error {
	start := time.Now()
	if err := datastore.Delete(cxt, key); err != nil {
		return err
	}
	total := time.Now().Sub(start)
	cxt.Infof("%d %s Single Delete: %v", count, key.Kind(), total)
	return nil
}

func putKinderSequential(w http.ResponseWriter, r *http.Request, kBuilder kinderBuilder) {
	outerStart := time.Now()
	cxt := appengine.NewContext(r)
	for i := 0; i < operationCount; i++ {
		entity := kBuilder.build()
		if err := putKinder(cxt, entity, i); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			cxt.Infof("%s", err.Error())
			return
		}
	}
	outerTotal := time.Now().Sub(outerStart)
	cxt.Infof("Few Kinds %d Puts: %v", operationCount, outerTotal)
}

func delKind(w http.ResponseWriter, r *http.Request, kind string) {
	cxt := appengine.NewContext(r)
	q := datastore.NewQuery(kind).KeysOnly().Limit(operationCount)
	keys, err := q.GetAll(cxt, make([]*empty, 0))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		cxt.Infof("%s", err.Error())
		return
	}
	for i, key := range keys {
		if err := delKey(cxt, key, i); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			cxt.Infof("%s", err.Error())
			return
		}
	}
	cxt.Infof("Sucessfully Deleted %s: %d", kind, len(keys))
}

func putKinderParallel(w http.ResponseWriter, r *http.Request, kBuilder kinderBuilder) {
	outerStart := time.Now()
	cxt := appengine.NewContext(r)
	complete := make(chan bool, operationCount)
	for i := 0; i < operationCount; i++ {
		go func(count int, entity kinder) {
			if err := putKinder(cxt, entity, count); err != nil {
				cxt.Infof("%s", err.Error())
			}
			complete <- true
		}(i, kBuilder.build())
	}
	outerTotal := time.Now().Sub(outerStart)
	cxt.Infof("Few Kinds %d Puts: %v", operationCount, outerTotal)
}
