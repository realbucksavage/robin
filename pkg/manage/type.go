package manage

import (
	"io"
	"sync"
)

type EventType int

const (
	Add EventType = iota
	Refresh
	Delete
)

type CertificateEvent struct {
	Type EventType
	Cert CertificateInfo
}

type CertificateInfo struct {
	HostName   string
	Cert       []byte
	PrivateKey []byte
	Origin     string
}

type CertEventBus interface {
	Subscribe(chan CertificateEvent)
	Emit(CertificateEvent)

	io.Closer
}

type defaultStore struct {
	subscribers []chan CertificateEvent
	mu          sync.RWMutex
}

func (d *defaultStore) Close() error {
	for _, c := range d.subscribers {
		close(c)
	}

	return nil
}

func (d *defaultStore) Subscribe(c chan CertificateEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.subscribers = append(d.subscribers, c)
}

func (d *defaultStore) Emit(event CertificateEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, s := range d.subscribers {
		s <- event
	}
}

func NewBus() CertEventBus {
	return &defaultStore{subscribers: make([]chan CertificateEvent, 0)}
}
