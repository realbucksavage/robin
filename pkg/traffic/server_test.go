package traffic

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/realbucksavage/robin/pkg/vhosts"
)

func TestTrafficHandler(t *testing.T) {
	backend := httptest.NewServer(dummyHandler())
	defer backend.Close()

	v, err := vhosts.TestingVault(t, backend)
	if err != nil {
		t.Fatalf("create dummy vault: %s", err)
	}

	proxy, err := NewProxy(v)
	if err != nil {
		t.Fatalf("create proxy: %s", err)
	}

	request := httptest.NewRequest("GET", "/", nil)
	request.Host = "localhost"

	responseRecorder := httptest.NewRecorder()

	proxy.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("got status code %d, expected %d", responseRecorder.Code, http.StatusOK)
	}
}

func TestTlsConfig(t *testing.T) {

	domains := []string{"localhost", "localhost.localdomain"}

	backend := httptest.NewServer(dummyHandler())
	defer backend.Close()

	v, err := vhosts.TestingVault(t, backend)
	if err != nil {
		t.Fatalf("create dummy vault: %s", err)
	}

	shutdown := make(chan bool)
	var wg sync.WaitGroup

	server := Server{
		Config: Config{
			BindAddr: ":8443",
		},
		VHostVault:   v,
		DoneFunc:     wg.Done,
		ShutdownChan: shutdown,
	}

	wg.Add(1)
	go server.Start()

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	for _, d := range domains {
		u := fmt.Sprintf("https://%s:8443/", d)
		get, err := client.Get(u)
		if err != nil {
			t.Fatal(err)
		}

		dnsNames := get.TLS.PeerCertificates[0].DNSNames
		t.Logf("Resolved DNS Names: %s", dnsNames)

		if dnsNames[0] != d {
			t.Fatalf("Expected %s; got %s in DNS name", d, dnsNames[0])
		}
	}

	if _, err = client.Get("https://127.0.0.1:8443/"); err == nil {
		t.Fatalf("expected \"unresolved certificate\" error")
	}

	close(shutdown)
	wg.Wait()

	t.Log("Test passed")
}

func dummyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`ok`))
	}
}
