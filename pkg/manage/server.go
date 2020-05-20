package manage

import (
	"context"
	"net/http"
	"time"

	"github.com/realbucksavage/robin/pkg/database"
	"github.com/realbucksavage/robin/pkg/log"
)

type Server struct {
	Config       Config
	ShutdownChan chan bool
	DoneFunc     func()
	Database     *database.Connection
	CertEventBus    CertEventBus
}

func (s *Server) Start() {
	bindAddr := s.Config.Bind
	if bindAddr == "" {
		bindAddr = defaultBindAddr
	}

	server := &http.Server{
		Addr:    bindAddr,
		Handler: newHandler(s.CertEventBus, s.Database),
	}
	log.L.Infof("Management interface active on %s", bindAddr)

	go gracefulShutdown(server, s.ShutdownChan, s.DoneFunc)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.L.Fatalf("management listen: %s", err)
	}
}

func gracefulShutdown(server *http.Server, shutdown chan bool, done func()) {
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.L.Errorf("management interface graceful shutdown: %s", err)
	}

	log.L.Info("Management interface closed.")
	done()
}
