package vhosts

import (
	"crypto/tls"
	"net/http/httputil"
)

type H struct {
	FQDN       string
	Origin     string
	PrivateKey []byte
	X509Cert   []byte
}

type Record struct {
	Backend *httputil.ReverseProxy
	Cert    *tls.Certificate
}

type Vault interface {
	Get(string) (Record, bool)
	Put(string, H) error
	Remove(string)
	Clear()
}
