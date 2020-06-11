package vhost

import (
	"context"
	"encoding/json"
	"errors"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var (
	errBadRouting = errors.New("bad route")
)

func MakeRouter(r *mux.Router, s Service) {
	e := MakeEndpoint(s)

	r.Methods("GET").Path("/").Handler(httptransport.NewServer(
		e.ListVhostsEndpoint,
		decodeListVhosts,
		encodeResponse,
	))
	r.Methods("DELETE").Path("/{id}").Handler(httptransport.NewServer(
		e.DeleteVhostEndpoint,
		decodeGetVhost,
		encodeResponse,
	))
	r.Methods("GET").Path("/{id}").Handler(httptransport.NewServer(
		e.GetVhostEndpoint,
		decodeGetVhost,
		encodeResponse,
	))
	r.Methods("POST").Path("/").Handler(httptransport.NewServer(
		e.CreateVhostEndpoint,
		decodePostVhost,
		encodeResponse,
	))
}

func decodeListVhosts(_ context.Context, r *http.Request) (interface{}, error) {
	// stub
	return struct{}{}, nil
}

func decodeGetVhost(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRouting
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return getVhostRequest{ID: uint(i)}, nil
}

func decodePostVhost(_ context.Context, r *http.Request) (interface{}, error) {
	var req createVhostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	type errorer interface {
		error() error
	}

	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}

	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
