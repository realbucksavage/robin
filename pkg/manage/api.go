package manage

import (
	"github.com/gorilla/mux"
	"github.com/realbucksavage/robin/pkg/database"
	"github.com/realbucksavage/robin/pkg/manage/api/vhost"
	"github.com/realbucksavage/robin/pkg/vhosts"
	"net/http"
)

type apiHandler struct {
	mux *mux.Router
}

func newHandler(store vhosts.Vault, conn *database.Connection, config AuthenticationConfig) (*apiHandler, error) {
	db, err := conn.Db()
	if err != nil {
		return nil, err
	}
	r := mux.NewRouter()
	r.Use(authenticationMiddleware(config))
	vh := r.PathPrefix("/api/vhosts").Subrouter()
	{
		vhost.MakeRouter(vh, vhost.NewService(db, store))
	}

	return &apiHandler{mux: r}, nil
}

func (a *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func authenticationMiddleware(config AuthenticationConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth, password, ok := r.BasicAuth()
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if auth != config.Username && password != config.Password {
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
