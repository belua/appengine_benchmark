package mercury

import (
	"appengine"
	"appengine/datastore"
	"net/http"
	"time"
)

const putCount = 200

type kinder interface {
	kind() string
}

type kinderBuilder interface {
	build() kinder
}

func init() {
	http.HandleFunc("/emptySequential", emptySequentialHandler)
	http.HandleFunc("/emptyParallel", emptyParallelHandler)
	http.HandleFunc("/delEmpty", delEmptyHandler)
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

func putKinderSequential(w http.ResponseWriter, r *http.Request, kBuilder kinderBuilder) {
	outerStart := time.Now()
	cxt := appengine.NewContext(r)
	for i := 0; i < putCount; i++ {
		entity := kBuilder.build()
		if err := putKinder(cxt, entity, i); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			cxt.Infof("%s", err.Error())
			return
		}
	}
	outerTotal := time.Now().Sub(outerStart)
	cxt.Infof("Few Kinds %d Puts: %v", putCount, outerTotal)
}

func putKinderParallel(w http.ResponseWriter, r *http.Request, kBuilder kinderBuilder) {
	outerStart := time.Now()
	cxt := appengine.NewContext(r)
	complete := make(chan bool, putCount)
	for i := 0; i < putCount; i++ {
		entity := kBuilder.build()
		go func(count int) {
			if err := putKinder(cxt, entity, count); err != nil {
				cxt.Infof("%s", err.Error())
			}
			complete <- true
		}(i)
	}
	outerTotal := time.Now().Sub(outerStart)
	cxt.Infof("Few Kinds %d Puts: %v", putCount, outerTotal)
}

func delKind(w http.ResponseWriter, r *http.Request, kind string) {
	cxt := appengine.NewContext(r)
	q := datastore.NewQuery(kind).KeysOnly().Limit(500)
	var keys []*datastore.Key
	var err error
	for keys == nil || len(keys) > 0 {
		keys, err = q.GetAll(cxt, make([]*empty, 0))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			cxt.Infof("%s", err.Error())
			return
		}
		cxt.Infof("Deleting %s: %d", kind, len(keys))
		if err := datastore.DeleteMulti(cxt, keys); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			cxt.Infof("%s", err.Error())
			return
		}
		cxt.Infof("Sucessfully Deleted %s: %d", kind, len(keys))
	}
}
