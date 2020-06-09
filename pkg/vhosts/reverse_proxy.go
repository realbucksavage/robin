package vhosts

import (
	"fmt"
	"github.com/realbucksavage/robin/pkg/log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func reverseProxy(backend string) (*httputil.ReverseProxy, error) {

	origin, err := url.Parse(backend)
	if err != nil {
		log.L.Errorf("Cannot initialize reverse proxy for %s: %s", backend, err)
		return nil, err
	}

	return &httputil.ReverseProxy{
		Director: defaultDirectorFunc(origin),
		ErrorLog: log.StdLogger,
		ErrorHandler: func(w http.ResponseWriter, e *http.Request, err error) {
			log.L.Warningf("Cannot reach backend %s: %s", backend, err)

			w.WriteHeader(http.StatusBadGateway)
			response := fmt.Sprintf(`<h1>Bad Gateway</h1><br><pre><code>%s</code></pre><hr>Robin v1.0`, err)
			w.Write([]byte(response))
		},
	}, nil
}

func defaultDirectorFunc(origin *url.URL) func(*http.Request) {
	return func(r *http.Request) {
		r.Header.Add("X-Forwarded-For", r.Host)
		r.Header.Add("X-Origin-Host", origin.Host)
		r.URL.Scheme = "http"

		r.URL.Host = origin.Host
	}
}
