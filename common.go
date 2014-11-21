package mercury

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
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
	start := time.Now.UnixNano()
	cxt := appengine.NewContext(r)
	e := &empty{}
	k := datastore.MakeIncompleteKey(cxt, e.kind(), nil)
	datastore.Put(cxt, k, e)
	total := time.Now.UnixNano() - start
	cxt.Logf("Few Kinds: %d", total)
}
