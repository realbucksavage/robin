package vhosts

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/url"
	"sync"

	"github.com/realbucksavage/robin/pkg/database"
	"github.com/realbucksavage/robin/pkg/log"
	"github.com/realbucksavage/robin/pkg/types"
)

type defaultVault struct {
	hostCerts map[string]Record
	mu        sync.RWMutex
}

func (v *defaultVault) Get(host string) (Record, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	h, ok := v.hostCerts[sanitize(host)]
	return h, ok
}

func (v *defaultVault) Put(host string, h H) error {
	host = sanitize(host)

	block, _ := pem.Decode(h.X509Cert)
	if block == nil {
		return fmt.Errorf("cannot parse block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("certificate parse error: %s", err)
	}

	if err := cert.VerifyHostname(host); err != nil {
		return fmt.Errorf("verify hostname: %s", err)
	}

	pair, err := tls.X509KeyPair(h.X509Cert, h.PrivateKey)
	if err != nil {
		return fmt.Errorf("invalid keypair: %v", err)
	}

	rp, err := reverseProxy(h.Origin)
	if err != nil {
		return err
	}

	v.mu.Lock()
	v.hostCerts[host] = Record{
		Backend: rp,
		Cert:    &pair,
	}
	v.mu.Unlock()

	log.L.Infof("Certificate stored for hostname %v", host)
	return nil
}

func (v *defaultVault) Remove(host string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	delete(v.hostCerts, sanitize(host))
	log.L.Infof("Certificate for host %s removed", host)
}

func (v *defaultVault) Clear() {
	v.mu.Lock()
	defer v.mu.Unlock()

	log.L.Info("Certificate defaultVault cleared")
	v.hostCerts = make(map[string]Record)
}

func NewVault(database *database.Connection) (Vault, error) {
	db, err := database.Db()
	if err != nil {
		return nil, err
	}

	var vhosts []types.Vhost
	if err := db.Preload("Cert").Find(&vhosts).Error; err != nil {
		return nil, err
	}

	s := defaultVault{
		hostCerts: make(map[string]Record),
	}
	for _, v := range vhosts {
		cert := v.Cert

		pair, err := tls.X509KeyPair(cert.X509, cert.RSAKey)
		if err != nil {
			log.L.Errorf("Invalid certificate and key pair for host %s", v.FQDN)
			continue
		}

		rp, err := reverseProxy(v.Origin)
		if err != nil {
			log.L.Errorf("Reverse proxy initialization failed (%s): %s", v.Origin, err)
			continue
		}

		host := sanitize(v.FQDN)
		s.hostCerts[host] = Record{
			Backend: rp,
			Cert:    &pair,
		}

		log.L.Debugf("Certificate for host %s added to the defaultVault", host)
	}

	return &s, nil
}

func sanitize(host string) string {
	u, err := url.Parse(host)
	if err != nil {
		log.L.Debugf("Host unchanged: %s", err)
	}

	newHost := u.Hostname()
	if newHost != "" {
		log.L.Debugf("Using host %s", newHost)
		host = newHost
	}
	return host
}
