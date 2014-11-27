package aebench

import (
	"appengine"
	"github.com/belua/httprouter"
	"net/http"
	"time"
)

func controlHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	start := time.Now()
	cxt := appengine.NewContext(r)
	total := time.Now().Sub(start)
	cxt.Infof("Control executed in %v", total)
}
