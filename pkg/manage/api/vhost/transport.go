package vhost

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	errBadRouting = errors.New("bad route")
)

func MakeRouter(r *mux.Router, s Service, mw endpoint.Middleware) {
	e := MakeEndpoint(s)

	r.Methods("GET").Path("/").Handler(httptransport.NewServer(
		mw(e.ListVhostsEndpoint),
		decodeListVhosts,
		httptransport.EncodeJSONResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	r.Methods("DELETE").Path("/{id}").Handler(httptransport.NewServer(
		mw(e.DeleteVhostEndpoint),
		decodeGetVhost,
		httptransport.EncodeJSONResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	r.Methods("GET").Path("/{id}").Handler(httptransport.NewServer(
		mw(e.GetVhostEndpoint),
		decodeGetVhost,
		httptransport.EncodeJSONResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	))
	r.Methods("POST").Path("/").Handler(httptransport.NewServer(
		mw(e.CreateVhostEndpoint),
		decodePostVhost,
		httptransport.EncodeJSONResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
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
