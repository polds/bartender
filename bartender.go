// bartender is an abstraction layer in front of
// go-cache designed specifically for Martini.
package bartender

import (
	"github.com/go-martini/martini"
	"github.com/pmylund/go-cache"
	"net/http"
	"time"
)

type Tab interface {
	Set(k string, x interface{}, d time.Duration)

	Get(k string) (interface{}, bool)
}

type tab struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	cache          *cache.Cache
}

func (t *tab) Set(k string, x interface{}, d time.Duration) {
	t.cache.Set(k, x, d)
}

func (t *tab) Get(k string) (interface{}, bool) {
	return t.cache.Get(k)
}

func NewTab(t *cache.Cache) martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		c.MapTo(&tab{res, req, t}, (*Tab)(nil))
	}
}

func OpenTab(defaultExpiration, cleanupInterval time.Duration) *cache.Cache {
	return cache.New(defaultExpiration, cleanupInterval)
}
