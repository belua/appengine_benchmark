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

type entityBuilder interface {
	build(appengine.Context) (*datastore.Key, interface{})
}

func init() {
	router := httprouter.New()
	router.GET("/control", controlHandler)
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

func putEntity(cxt appengine.Context, key *datastore.Key, entity interface{}, count int) error {
	start := time.Now()
	if _, err := datastore.Put(cxt, key, entity); err != nil {
		return err
	}
	total := time.Now().Sub(start)
	cxt.Infof("%d %s Single Put: %v", count, key.Kind(), total)
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

func putEntitySequential(w http.ResponseWriter, r *http.Request, eBuilder entityBuilder) {
	outerStart := time.Now()
	cxt := appengine.NewContext(r)
	for i := 0; i < operationCount; i++ {
		key, entity := eBuilder.build(cxt)
		if err := putEntity(cxt, key, entity, i); err != nil {
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
	// url := "http://localhost:8080" + params.ByName("url")
	url := "http://fiery-diorama-771.appspot.com" + params.ByName("url")
	complete := make(chan bool, size)
	for i := 0; i < size; i++ {
		go func(idx int) {
			start := time.Now()
			client := urlfetch.Client(cxt)
			_, err := client.Get(url)
			if err != nil {
				cxt.Infof("URL Fetch error: %s", err.Error())
			}
			total := time.Now().Sub(start)
			cxt.Infof("URL fetch %d complete in %v", idx, total)
			complete <- true
		}(i)
	}
	for i := 0; i < size; i++ {
		<-complete
	}
	outerTotal := time.Now().Sub(outerStart)
	cxt.Infof("Loaded %s %d times in %v", url, size, outerTotal)
}
