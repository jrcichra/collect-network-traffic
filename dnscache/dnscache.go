package dnscache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

//DNSCache - Manages state of our DNS cache
type DNSCache struct {
	cache *cache.Cache
}

//Start - starts a new DNS cache
func (c *DNSCache) Start() {
	c.cache = cache.New(10*time.Second, 20*time.Second)
}

//Set - Sets a new DNS record
func (c *DNSCache) Set(k, v string) {
	c.cache.Set(k, v, cache.DefaultExpiration)
}

//Get - Get a new DNS record
func (c *DNSCache) Get(k string) (interface{}, bool) {
	return c.cache.Get(k)
}
