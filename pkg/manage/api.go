package manage

import (
	"encoding/json"
	"github.com/realbucksavage/robin/pkg/vhosts"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/realbucksavage/robin/pkg/database"
	"github.com/realbucksavage/robin/pkg/log"
	"github.com/realbucksavage/robin/pkg/types"
	"gopkg.in/go-playground/validator.v8"
)

type apiHandler struct {
	mux   *mux.Router
	conn  *database.Connection
	store vhosts.Vault
}

func newHandler(store vhosts.Vault, conn *database.Connection) *apiHandler {
	r := mux.NewRouter()

	r.HandleFunc("/api/hosts", createHost(store, conn)).Methods("POST")

	return &apiHandler{mux: r, store: store}
}

func (a *apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}

func createHost(s vhosts.Vault, conn *database.Connection) http.HandlerFunc {

	type request struct {
		FQDN     string `json:"fqdn" validate:"required"`
		Nickname string `json:"nick"`
		Key      string `json:"rsa" validate:"required"`
		Cert     string `json:"cert" validate:"reqired"`
		Origin   string `json:"origin" validate:"required"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.L.Errorf("read body: %s", err)
			return
		}

		v := validator.New(&validator.Config{})

		var rq request
		if err := json.Unmarshal(bytes, &rq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := v.Struct(rq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if rq.Nickname == "" {
			rq.Nickname = rq.FQDN
		}

		db, err := conn.Db()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		keyBytes := []byte(rq.Key)
		certBytes := []byte(rq.Cert)

		tx := db.Begin()
		if err := tx.Save(&types.Vhost{
			FQDN:   rq.FQDN,
			Origin: rq.Origin,
			Cert: types.Certificate{
				RSAKey:  keyBytes,
				X509:    certBytes,
				CAChain: nil,
			},
		}).Error; err != nil {
			log.L.Errorf("save certificate: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()

			return
		}

		if err = s.Put(rq.FQDN, vhosts.H{
			FQDN:       rq.FQDN,
			Origin:     rq.Origin,
			PrivateKey: keyBytes,
			X509Cert:   certBytes,
		}); err != nil {
			log.L.Errorf("update certificate store: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			tx.Rollback()

			return
		}

		tx.Commit()

		w.WriteHeader(http.StatusCreated)
		log.L.Infof("New certificate saved for host %s", rq.FQDN)
	}
}
