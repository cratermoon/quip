package svc

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type TransientStorageService struct {
	sync.RWMutex
	store map[string]time.Time
}

var tss TransientStorageService

type TransientStorable interface {
	Value() string
}

func add(key string) {
	tss.Lock()
	tss.store[key] = time.Now()	
	tss.Unlock()
}

func remove(key string) {
	tss.Lock()
	delete(tss.store, key)
	tss.Unlock()
}

func reap() {
	for {
		time.Sleep(time.Second)
		tss.RLock()
		log.Printf("Reaping store at %s:  %d entries\n", time.Now().Format(time.StampMilli), len(tss.store))
		for key, value := range tss.store {
			if time.Since(value) > time.Duration(3 * time.Second) {
				log.Printf("Deleting, %s entries\n", key)
				delete(tss.store, key)
			}
		}
		tss.RUnlock()
	}
}

// MakeStorageMiddleware decorates a call to store certain values
func MakeStorageMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			resp, error := next(ctx, request)
			if error != nil {
				return resp, error
			}
			r, ok := resp.(TransientStorable)
			log.Printf("Adding? %t\n", ok)
			if ok {
				add(r.Value())
			}
			return resp, error
		}
	}

}

// MakeLookupMiddleware decorates a call to return certain values
func MakeLookupMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			r, ok := request.(TransientStorable)
			if ok {
				remove(r.Value())
			}
			return next(ctx, request)
		}
	}

}

func init() {
	// create the storage
	tss = TransientStorageService{}
	tss.store = make(map[string]time.Time)
	// start the go routine to reap old values
	go reap()
}
