package aebench

import (
	"appengine"
	"appengine/datastore"
	"github.com/belua/httprouter"
	"math/rand"
	"net/http"
	"time"
)

// Entity with single indexable field

type OneIndex struct {
	One int64
}

type OneIndexBuilder struct {
	r *rand.Rand
}

func (b *OneIndexBuilder) build(cxt appengine.Context) (*datastore.Key, interface{}) {
	return datastore.NewIncompleteKey(cxt, "OneIndex", nil), OneIndex{b.r.Int63()}
}

func oneIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &OneIndexBuilder{rand.New(rand.NewSource(time.Now().UnixNano()))}
	putEntities(w, r, b)
}

// Entity with two indexable fields

type TwoIndex struct {
	One int64
	Two int64
}

type TwoIndexBuilder struct {
	r *rand.Rand
}

func (b *TwoIndexBuilder) build(cxt appengine.Context) (*datastore.Key, interface{}) {
	return datastore.NewIncompleteKey(cxt, "TwoIndex", nil), TwoIndex{b.r.Int63(), b.r.Int63()}
}

func twoIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &TwoIndexBuilder{rand.New(rand.NewSource(time.Now().UnixNano()))}
	putEntities(w, r, b)
}

// Entity with three indexable fields

type ThreeIndex struct {
	One   int64
	Two   int64
	Three int64
}

type ThreeIndexBuilder struct {
	r *rand.Rand
}

func (b *ThreeIndexBuilder) build(cxt appengine.Context) (*datastore.Key, interface{}) {
	return datastore.NewIncompleteKey(cxt, "ThreeIndex", nil), ThreeIndex{b.r.Int63(), b.r.Int63(), b.r.Int63()}
}

func threeIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &ThreeIndexBuilder{rand.New(rand.NewSource(time.Now().UnixNano()))}
	putEntities(w, r, b)
}

// Entity with four indexable fields

type FourIndex struct {
	One   int64
	Two   int64
	Three int64
	Four  int64
}

type FourIndexBuilder struct {
	r *rand.Rand
}

func (b *FourIndexBuilder) build(cxt appengine.Context) (*datastore.Key, interface{}) {
	return datastore.NewIncompleteKey(cxt, "FourIndex", nil), FourIndex{b.r.Int63(), b.r.Int63(), b.r.Int63(), b.r.Int63()}
}

func fourIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &FourIndexBuilder{rand.New(rand.NewSource(time.Now().UnixNano()))}
	putEntities(w, r, b)
}
