package traffic

import (
	"context"
	"crypto/tls"
	"github.com/realbucksavage/robin/pkg/vhosts"
	"net/http"
	"time"

	"github.com/realbucksavage/robin/pkg/log"
)

type Server struct {
	Config       Config
	ShutdownChan chan bool
	DoneFunc     func()
	VHostVault   vhosts.Vault
}

func (s *Server) Start() {
	bindAddr := s.Config.BindAddr
	if bindAddr == "" {
		bindAddr = defaultBindAddr
	}

	tlsConfig := &tls.Config{
		GetCertificate: getCertificateFunc(s.VHostVault),
	}

	listener, err := tls.Listen("tcp", bindAddr, tlsConfig)
	if err != nil {
		log.L.Fatalf("create traffic listener: %s", err)
	}

	handler, err := NewProxy(s.VHostVault)
	if err != nil {
		log.L.Fatalf("create proxy: %s", err)
	}

	server := &http.Server{
		TLSConfig: tlsConfig,
		Handler:   handler,
		ErrorLog:  log.StdLogger,
	}
	go gracefulShutdown(server, s.ShutdownChan, s.DoneFunc)

	log.L.Infof("Listening for HTTPs traffic on %s", bindAddr)
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		log.L.Fatalf("traffic listen: %s", err)
	}
}

func gracefulShutdown(server *http.Server, shutdown chan bool, done func()) {
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.L.Errorf("Traffic interface failed to shutdown gracefully: %s", err)
	}

	log.L.Info("Traffic interface closed.")
	done()
}
