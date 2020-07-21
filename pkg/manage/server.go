package manage

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/realbucksavage/robin/pkg/vhosts"

	"github.com/realbucksavage/robin/pkg/database"
	"github.com/realbucksavage/robin/pkg/log"
)

type Server struct {
	Config       Config
	ShutdownChan chan bool
	DoneFunc     func()
	Database     *database.Connection
	VHostVault   vhosts.Vault
}

func (s *Server) Start() {
	bindAddr := s.Config.Bind
	if bindAddr == "" {
		bindAddr = defaultBindAddr
	}

	auth := s.Config.Auth
	if auth.Username == "" {
		log.L.Infof("Using the default username (%s) for management API", defaultUsername)
		auth.Username = defaultUsername
	}

	if auth.Password == "" {
		p := randomPassword(10)

		log.L.Infof("Using a generated password (%s) for management API", p)
		auth.Password = p
	}

	handler, err := newHandler(s.VHostVault, s.Database, s.Config.Auth)
	if err != nil {
		log.L.Fatalf("create handler: %s", err)
	}
	server := &http.Server{
		Addr:    bindAddr,
		Handler: handler,
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

func randomPassword(l int) string {
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, l)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
