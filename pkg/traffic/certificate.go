package traffic

import (
	"crypto/tls"
	"net/url"
	"sync"

	"github.com/realbucksavage/robin/pkg/database"
	"github.com/realbucksavage/robin/pkg/log"
	"github.com/realbucksavage/robin/pkg/manage"
	"github.com/realbucksavage/robin/pkg/types"
)

type certificateCache struct {
	mu    sync.RWMutex
	certs []manage.CertificateInfo
}

func getCertificateFunc(
	bus manage.CertEventBus,
	conn *database.Connection,
) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {

	db, err := conn.Db()
	if err != nil {
		return errorFunc(err)
	}

	cache := &certificateCache{
		certs: make([]manage.CertificateInfo, 0),
	}

	var hosts []types.Host
	if err := db.Find(&hosts).Error; err != nil {
		return errorFunc(err)
	}

	for _, h := range hosts {
		cache.certs = append(cache.certs, manage.CertificateInfo{
			HostName:   h.FQDN,
			Cert:       h.SSLCertificate,
			PrivateKey: h.RSAKey,
		})
	}

	sub := make(chan manage.CertificateEvent)
	go subscribe(cache, sub)

	bus.Subscribe(sub)

	return func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
		cache.mu.RLock()
		defer cache.mu.RUnlock()

		for _, c := range cache.certs {
			origin, err := url.Parse(c.HostName)
			if err != nil {
				log.L.Errorf("URL %s can't be parsed: %s", c.HostName, err)
				continue
			}

			if info.ServerName == origin.Host {
				cert, err := tls.X509KeyPair(c.Cert, c.PrivateKey)
				if err != nil {
					return nil, err
				}

				return &cert, nil
			}
		}

		return nil, nil
	}
}

func subscribe(cache *certificateCache, sub chan manage.CertificateEvent) {
	for event := range sub {
		cache.mu.Lock()

		switch event.Type {
		case manage.Add:
			cache.certs = append(cache.certs, event.Cert)
			log.L.Infof("Added certificate for %s to the cert store.", event.Cert.HostName)
			break
		case manage.Delete:
			newCerts := make([]manage.CertificateInfo, 0)
			for _, c := range cache.certs {
				if c.HostName != event.Cert.HostName {
					newCerts = append(newCerts, event.Cert)
				}
			}

			cache.certs = newCerts
			log.L.Infof("Removed certificate for host %s from the proxyServers", event.Cert.HostName)
		}

		cache.mu.Unlock()
	}
}

func errorFunc(err error) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		log.L.Errorf("cannot initialize initial pool of certificates: %s", err)
		return nil, err
	}
}
