package traffic

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/realbucksavage/robin/pkg/database"
	"github.com/realbucksavage/robin/pkg/log"
	"github.com/realbucksavage/robin/pkg/manage"
	"github.com/realbucksavage/robin/pkg/types"
)

type reverseProxyCache struct {
	mu           sync.RWMutex
	proxyServers map[string]*httputil.ReverseProxy
}

func NewProxy(bus manage.CertEventBus, conn *database.Connection) (http.Handler, error) {

	db, err := conn.Db()
	if err != nil {
		return nil, err
	}

	cache := &reverseProxyCache{proxyServers: make(map[string]*httputil.ReverseProxy)}

	var hosts []types.Host
	if err := db.Find(&hosts).Error; err != nil {
		return nil, err
	}

	for _, h := range hosts {
		addReverseProxyToCache(h.FQDN, h.Origin, cache)
	}

	subChan := make(chan manage.CertificateEvent)
	bus.Subscribe(subChan)

	go hostCacheSubscriber(cache, subChan)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cache.mu.RLock()
		defer cache.mu.RUnlock()

		hostWithoutPort := strings.Split(r.Host, ":")[0]
		if rp, ok := cache.proxyServers[hostWithoutPort]; ok && rp != nil {
			rp.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write(statusText(http.StatusServiceUnavailable))
		}
	}), nil
}

func hostCacheSubscriber(cache *reverseProxyCache, subChan chan manage.CertificateEvent) {
	for event := range subChan {
		cache.mu.Lock()

		switch event.Type {
		case manage.Add:
			addReverseProxyToCache(event.Cert.HostName, event.Cert.Origin, cache)
			log.L.Infof("Host %s (proxies %S) added", event.Cert.HostName, event.Cert.Origin)
			break
		case manage.Delete:
			log.L.Notice("TODO: Implement deleting of hosts from proxy")
		}

		cache.mu.Unlock()
	}
}

func addReverseProxyToCache(hostname, backend string, cache *reverseProxyCache) {
	origin, err := url.Parse(hostname)
	if err != nil {
		log.L.Errorf("Host %s can't be parsed: %s", hostname, err)
		return
	}

	proxy := &httputil.ReverseProxy{
		Director: defaultDirectorFunc(backend),
		ErrorLog: log.StdLogger,
	}

	cache.proxyServers[origin.Hostname()] = proxy
}

func defaultDirectorFunc(backend string) func(*http.Request) {
	origin, err := url.Parse(backend)
	if err != nil {
		log.L.Errorf("Cannot initialize reverse proxy for %s: %s", backend, err)
		return nil
	}

	return func(r *http.Request) {
		r.Header.Add("X-Forwarded-For", r.Host)
		r.Header.Add("X-Origin-Host", origin.Host)
		r.URL.Scheme = "http"

		r.URL.Host = origin.Host
	}
}

func statusText(status int) []byte {
	t := fmt.Sprintf("<h1>%s</h1><hr>Status code %d", http.StatusText(status), status)
	return []byte(t)
}
