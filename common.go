package aebench

import (
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
	"github.com/belua/httprouter"
	"net/http"
	"strconv"
	"time"
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
	router.GET("/oneIndex", oneIndexHandler)
	router.GET("/twoIndex", twoIndexHandler)
	router.GET("/threeIndex", threeIndexHandler)
	router.GET("/fourIndex", fourIndexHandler)
	router.GET("/monoIndex", monoIndexHandler)
	router.GET("/del/:kind", delHandler)
	router.GET("/clear/:kind", clearHandler)
	router.GET("/load/:size/*url", loadHandler)
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

func delHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	kind := params.ByName("kind")
	outerStart := time.Now()
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
	outerTotal := time.Now().Sub(outerStart)
	cxt.Infof("Sucessfully Deleted %s: %d in %v", kind, len(keys), outerTotal)
}

func clearHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	kind := params.ByName("kind")
	outerStart := time.Now()
	cxt := appengine.NewContext(r)
	q := datastore.NewQuery(kind).KeysOnly().Limit(500)
	for {
		keys, err := q.GetAll(cxt, make([]*empty, 0))
		if len(keys) == 0 {
			break
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			cxt.Infof("%s", err.Error())
			return
		}
		if err := datastore.DeleteMulti(cxt, keys); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			cxt.Infof("%s", err.Error())
			return
		}
	}
	outerTotal := time.Now().Sub(outerStart)
	cxt.Infof("Sucessfully Cleared %s in %v", kind, outerTotal)
}

func loadHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	outerStart := time.Now()
	cxt := appengine.NewContext(r)
	size, err := strconv.Atoi(params.ByName("size"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		cxt.Infof("%s", err.Error())
		return
	}
	url := params.ByName("url")
	complete := make(chan bool, size)
	for i := 0; i < size; i++ {
		go func() {
			client := urlfetch.Client(cxt)
			// _, err := client.Get("http://fiery-diorama-777.appspot.com/"+url)
			_, err := client.Get("http://localhost:8080/" + url)
			if err != nil {
				cxt.Infof("URL Fetch error: %s", err.Error())
			}
			complete <- true
		}()
	}
	for i := 0; i < size; i++ {
		<-complete
	}
	outerTotal := time.Now().Sub(outerStart)
	cxt.Infof("Loaded %s %d times in %v", url, size, outerTotal)
}
