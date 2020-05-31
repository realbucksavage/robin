package vhosts

import (
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
