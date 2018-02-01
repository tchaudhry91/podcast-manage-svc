package service

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

// Endpoints is a struct which contains a full list of endpoints for the PodcastManageService
type Endpoints struct {
	CreateUserEndpoint           endpoint.Endpoint
	GetUserEndpoint              endpoint.Endpoint
	GetPodcastDetailsEndpoint    endpoint.Endpoint
	SubscribeEndpoint            endpoint.Endpoint
	GetUserSubscriptionsEndpoint endpoint.Endpoint
	GetTokenEndpoint             endpoint.Endpoint
}

// MakeServerEndpoints returns a struct containing all the endpoints for a PodcastManageService
func MakeServerEndpoints(svc PodcastManageService) Endpoints {
	return Endpoints{
		CreateUserEndpoint: MakeCreateUserEndpoint(svc),
	}
}

// MakeCreateUserEndpoint returns a CreateUserEndpoint via the passed service
func MakeCreateUserEndpoint(svc PodcastManageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createUserRequest)
		e := svc.CreateUser(ctx, req.EmailID, req.Password)
		if e != nil {
			return createUserResponse{e.Error()}, nil
		}
		return createUserResponse{Err: ""}, nil
	}
}

type createUserRequest struct {
	EmailID  string `json:"email_id"`
	Password string `json:"password"`
}

type createUserResponse struct {
	Err string `json:"err,omitempty"`
}
