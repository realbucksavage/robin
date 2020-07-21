package vhosts

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

var (
	errNotImplemented = errors.New("not implemented")
)

type certPair struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

type dummyVault struct {
	hostCerts map[string]Record
	t         *testing.T
	s         *httptest.Server
}

func (d *dummyVault) Get(in string) (Record, bool) {
	h, ok := d.hostCerts[sanitize(in)]
	return h, ok
}

func (d *dummyVault) Put(_ string, _ H) error {
	return errNotImplemented
}

func (d *dummyVault) Remove(_ string) {
}

func (d *dummyVault) Clear() {
}

func TestingVault(t *testing.T, s *httptest.Server) (Vault, error) {
	raw, err := ioutil.ReadFile("../testdata/certificates.json")
	if err != nil {
		return nil, err
	}

	hosts := make(map[string]certPair)
	if err := json.Unmarshal(raw, &hosts); err != nil {
		return nil, err
	}

	vault := &dummyVault{hostCerts: make(map[string]Record), t: t, s: s}
	for h, p := range hosts {

		t.Logf("Adding host %s: %s", h, p.Cert)

		pair, err := tls.X509KeyPair([]byte(p.Cert), []byte(p.Key))
		if err != nil {
			return nil, err
		}

		host := sanitize(h)
		backend := "http://" + s.Listener.Addr().String()
		rp, err := reverseProxy(backend)
		if err != nil {
			return nil, err
		}
		vault.hostCerts[host] = Record{
			Cert:    &pair,
			Backend: rp,
		}
	}

	return vault, nil
}
