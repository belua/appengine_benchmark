package aebench

import (
	"github.com/belua/httprouter"
	"math/rand"
	"net/http"
	"time"
)

// Entity with single indexable field

type OneIndex struct {
	One int64
}

func (i *OneIndex) kind() string {
	return "OneIndex"
}

type OneIndexBuilder struct {
	r *rand.Rand
}

func (b *OneIndexBuilder) build() kinder {
	return &OneIndex{b.r.Int63()}
}

func oneIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &OneIndexBuilder{rand.New(rand.NewSource(time.Now().UnixNano()))}
	putKinderSequential(w, r, b)
}

// Entity with two indexable fields

type TwoIndex struct {
	One int64
	Two int64
}

func (i *TwoIndex) kind() string {
	return "TwoIndex"
}

type TwoIndexBuilder struct {
	r *rand.Rand
}

func (b *TwoIndexBuilder) build() kinder {
	return &TwoIndex{b.r.Int63(), b.r.Int63()}
}

func twoIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &TwoIndexBuilder{rand.New(rand.NewSource(time.Now().UnixNano()))}
	putKinderSequential(w, r, b)
}

// Entity with three indexable fields

type ThreeIndex struct {
	One   int64
	Two   int64
	Three int64
}

func (i *ThreeIndex) kind() string {
	return "ThreeIndex"
}

type ThreeIndexBuilder struct {
	r *rand.Rand
}

func (b *ThreeIndexBuilder) build() kinder {
	return &ThreeIndex{b.r.Int63(), b.r.Int63(), b.r.Int63()}
}

func threeIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &ThreeIndexBuilder{rand.New(rand.NewSource(time.Now().UnixNano()))}
	putKinderSequential(w, r, b)
}

// Entity with four indexable fields

type FourIndex struct {
	One   int64
	Two   int64
	Three int64
	Four  int64
}

func (i *FourIndex) kind() string {
	return "FourIndex"
}

type FourIndexBuilder struct {
	r *rand.Rand
}

func (b *FourIndexBuilder) build() kinder {
	return &FourIndex{b.r.Int63(), b.r.Int63(), b.r.Int63(), b.r.Int63()}
}

func fourIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b := &FourIndexBuilder{rand.New(rand.NewSource(time.Now().UnixNano()))}
	putKinderSequential(w, r, b)
}
