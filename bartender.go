// bartender is an abstraction layer in front of
// go-cache designed specifically for Martini.
package bartender

import (
	"github.com/go-martini/martini"
	"github.com/pmylund/go-cache"
	"net/http"
	"time"
)

type Item struct {
	Object     interface{}
	Expiration *time.Time
}

type Tab interface {
	Set(k string, x interface{}, d time.Duration)
	Get(k string) (interface{}, bool)
	DeleteExpired()
	Items() map[string]*cache.Item
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

// DeleteExpired is the abstraction of cache.DeleteExpired() used
// to force a manual purge of already expired items in the cache
// before the janitor gets to it.
func (t *tab) DeleteExpired() {
	t.cache.DeleteExpired()
}

// Items is the abstraction of cache.Items() which
// returns all items in the cache including the potential of
// expired items. Please see their documentation for more info.
func (t *tab) Items() map[string]*cache.Item {
	return t.cache.Items()
}

func NewTab(t *cache.Cache) martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		c.MapTo(&tab{res, req, t}, (*Tab)(nil))
	}
}

func OpenTab(defaultExpiration, cleanupInterval time.Duration) *cache.Cache {
	return cache.New(defaultExpiration, cleanupInterval)
}
