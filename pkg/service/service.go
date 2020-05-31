package service

import (
	"github.com/realbucksavage/robin/pkg/vhosts"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/realbucksavage/robin/pkg/database"
	"github.com/realbucksavage/robin/pkg/log"
	"github.com/realbucksavage/robin/pkg/manage"
	"github.com/realbucksavage/robin/pkg/traffic"
	"github.com/realbucksavage/robin/pkg/types"
)

func cmdService(config Config) {

	db, err := database.NewConnection(config.Database)
	if err != nil {
		log.L.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.L.Errorf("database graceful shutdown: %s", err)
		}
	}()

	if err := db.Migrate(
		&types.Vhost{},
		&types.Certificate{},
	); err != nil {
		log.L.Fatal(err)
	}

	vhostVault, err := vhosts.NewVault(db)
	if err != nil {
		log.L.Fatal(err)
	}

	shutdown := make(chan bool)
	var wg sync.WaitGroup

	trafficServer := &traffic.Server{
		Config:       config.Traffic,
		ShutdownChan: shutdown,
		DoneFunc:     wg.Done,
		VHostVault:   vhostVault,
	}
	wg.Add(1)
	go trafficServer.Start()

	managementServer := &manage.Server{
		Config:       config.Management,
		ShutdownChan: shutdown,
		DoneFunc:     wg.Done,
		Database:     db,
		VHostVault:   vhostVault,
	}
	wg.Add(1)
	go managementServer.Start()

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-interrupt

	log.L.Info("Shutting down...")
	close(shutdown)

	wg.Wait()
	log.L.Info("A long night in Gotham is finally over.")
}
