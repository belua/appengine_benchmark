package aebench

import (
	"appengine"
	"appengine/datastore"
	"net/http"
	"time"
	"github.com/belua/httprouter"
)

const operationCount = 20

type kinder interface {
	kind() string
}

type kinderBuilder interface {
	build() kinder
}

func init() {
	router := httprouter.New()
	router.GET("/empty", emptyHandler)
	router.GET("/emptyDel", emptyDelHandler)
	router.GET("/oneIndex", oneIndexHandler)
	router.GET("/oneIndexDel", oneIndexDelHandler)
	router.GET("/twoIndex", twoIndexHandler)
	router.GET("/twoIndexDel", twoIndexDelHandler)
	router.GET("/threeIndex", threeIndexHandler)
	router.GET("/threeIndexDel", threeIndexDelHandler)
	router.GET("/fourIndex", fourIndexHandler)
	router.GET("/fourIndexDel", fourIndexDelHandler)
	router.GET("/monoIndex", monoIndexHandler)
	router.GET("/monoIndexDel", monoIndexDelHandler)
	http.Handle("/", router)
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
