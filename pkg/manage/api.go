package manage

import (
	"net/http"

	"github.com/go-kit/kit/auth/basic"
	"github.com/gorilla/mux"
	"github.com/realbucksavage/robin/pkg/database"
	"github.com/realbucksavage/robin/pkg/manage/api/vhost"
	"github.com/realbucksavage/robin/pkg/vhosts"
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

	mw := basic.AuthMiddleware(config.Username, config.Password, "robin")
	vh := r.PathPrefix("/api/vhosts").Subrouter()
	{
		vhost.MakeRouter(vh, vhost.NewService(db, store), mw)
	}

	return &apiHandler{mux: r}, nil
}

func (a *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
