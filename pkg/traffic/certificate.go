package traffic

import (
	"crypto/tls"
	"fmt"
	"github.com/realbucksavage/robin/pkg/log"
	"github.com/realbucksavage/robin/pkg/vhosts"
)

func getCertificateFunc(store vhosts.Vault) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {

	if store == nil {
		return func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return nil, fmt.Errorf("certificate store is nil")
		}
	}

	return func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {

		vh, ok := store.Get(info.ServerName)
		if !ok {
			log.L.Warningf("No matching certificate found for server %s", info.ServerName)
			return nil, nil
		}

		return vh.Cert, nil
	}
}
