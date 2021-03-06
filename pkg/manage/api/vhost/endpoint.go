package vhost

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/realbucksavage/robin/pkg/types"
	"gopkg.in/go-playground/validator.v8"
)

type Endpoints struct {
	ListVhostsEndpoint  endpoint.Endpoint
	GetVhostEndpoint    endpoint.Endpoint
	CreateVhostEndpoint endpoint.Endpoint
	DeleteVhostEndpoint endpoint.Endpoint
}

type getVhostRequest struct {
	ID uint
}

type createVhostRequest struct {
	FQDN   string `json:"fqdn" validate:"required"`
	Key    string `json:"rsa" validate:"required"`
	Cert   string `json:"cert" validate:"required"`
	Origin string `json:"origin" validate:"required"`
}

func MakeEndpoint(s Service) Endpoints {

	return Endpoints{
		ListVhostsEndpoint:  makeListVhostsEndpoint(s),
		GetVhostEndpoint:    makeGetVhostEndpoint(s),
		CreateVhostEndpoint: makePostVhostEndpoint(s),
		DeleteVhostEndpoint: makeDeleteVhostEndpoiint(s),
	}
}

func makeListVhostsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if v, err := s.ListVhosts(ctx); err != nil {
			return nil, err
		} else {
			return v, nil
		}
	}
}

func makeDeleteVhostEndpoiint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getVhostRequest)
		if err := s.DeleteVhost(ctx, req.ID); err != nil {
			return nil, err
		} else {
			return map[string]interface{}{
				"success": true,
			}, nil
		}
	}
}

func makeGetVhostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getVhostRequest)
		if v, err := s.GetVhost(ctx, req.ID); err != nil {
			return nil, err
		} else {
			return v, nil
		}
	}
}

func makePostVhostEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createVhostRequest)

		v := validator.New(&validator.Config{})
		if err := v.Struct(v); err != nil {
			return nil, err
		}

		vhost, err := s.PostVhost(ctx, types.Vhost{
			FQDN:   req.FQDN,
			Origin: req.Origin,
			Cert: types.Certificate{
				RSAKey:  []byte(req.Key),
				X509:    []byte(req.Cert),
				CAChain: make([]byte, 0),
			},
		})
		if err != nil {
			return nil, err
		}

		return vhost, nil
	}
}
