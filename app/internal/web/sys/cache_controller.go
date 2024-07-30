package internal_sys

import (
	"net/http"
	internalcache "std-library/app/internal/cache"
	"std-library/app/internal/web/http"
	"std-library/app/web/errors"
	"std-library/json"
	"std-library/nets"
	"strings"
)

type CacheController struct {
	caches        map[string]*internalcache.CacheImpl
	accessControl *internal_http.IPv4AccessControl
}

func NewCacheController(caches map[string]*internalcache.CacheImpl) *CacheController {
	return &CacheController{
		caches:        caches,
		accessControl: &internal_http.IPv4AccessControl{},
	}
}

func (c *CacheController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.accessControl.Validate(nets.IP(r).String())
	if err != nil {
		errors.Forbidden("access denied", "IP_ACCESS_DENIED")
	}

	if r.Method == http.MethodGet && r.URL.Path == "/_sys/cache" {
		views := c.list()
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(json.Stringify(views))
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/_sys/cache/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		errors.NotFound("not found")
	}

	name := parts[0]
	key := parts[1]
	cache := c.cache(name)
	if r.Method == http.MethodDelete {
		cache.Evict(r.Context(), key)
		w.WriteHeader(200)

		w.Write([]byte("cache evicted, name=" + name + ", key=" + key))
		return
	}

	if r.Method == http.MethodGet {
		value := make(map[string]interface{})
		ok := cache.GetByKey(r.Context(), key, &value)
		if !ok {
			errors.NotFoundError(404, "cache key not found, name="+name+", key="+key)
		}
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(json.Stringify(value))
		return
	}
	errors.NotFound("not found")
}

func (c *CacheController) list() []CacheView {
	var views []CacheView
	for _, cache := range c.caches {
		views = append(views, view(cache))
	}
	return views
}

func (c *CacheController) cache(name string) *internalcache.CacheImpl {
	impl, ok := c.caches[name]
	if !ok {
		errors.NotFoundError(404, "cache not found: "+name)
	}
	return impl
}

func view(cache *internalcache.CacheImpl) CacheView {
	return CacheView{
		Name:     cache.Name(),
		Type:     cache.GetTypeName(),
		Duration: cache.Expiration.String(),
	}
}

type CacheView struct {
	Name     string
	Type     string
	Duration string
}
