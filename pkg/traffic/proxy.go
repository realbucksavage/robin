package traffic

import (
	"fmt"
	"github.com/realbucksavage/robin/pkg/vhosts"
	"net/http"
	"strings"
)

func NewProxy(store vhosts.Vault) (http.Handler, error) {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hostWithoutPort := strings.Split(r.Host, ":")[0]
		vhost, ok := store.Get(hostWithoutPort)
		if !ok {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write(statusText(http.StatusServiceUnavailable))
			return
		}

		vhost.Backend.ServeHTTP(w, r)
	}), nil
}


func statusText(status int) []byte {
	t := fmt.Sprintf("<h1>%s</h1><hr>Status code %d", http.StatusText(status), status)
	return []byte(t)
}
