package manage

import (
	"github.com/realbucksavage/robin/pkg/manage/api/vhost"
	"github.com/realbucksavage/robin/pkg/vhosts"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/realbucksavage/robin/pkg/database"
)

type apiHandler struct {
	mux *mux.Router
}

func newHandler(store vhosts.Vault, conn *database.Connection) (*apiHandler, error) {
	db, err := conn.Db()
	if err != nil {
		return nil, err
	}
	r := mux.NewRouter()
	vh := r.PathPrefix("/api/vhosts").Subrouter()
	{
		vhost.MakeRouter(vh, vhost.NewService(db, store))
	}

	return &apiHandler{mux: r}, nil
}

func (a *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
