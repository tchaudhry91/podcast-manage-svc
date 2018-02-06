package service

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/tchaudhry91/podcast-manage-svc/podcastmg"
)

// Endpoints is a struct which contains a full list of endpoints for the PodcastManageService
type Endpoints struct {
	CreateUserEndpoint             endpoint.Endpoint
	GetUserEndpoint                endpoint.Endpoint
	GetPodcastDetailsEndpoint      endpoint.Endpoint
	SubscribeEndpoint              endpoint.Endpoint
	GetUserSubscriptionsEndpoint   endpoint.Endpoint
	GetSubscriptionDetailsEndpoint endpoint.Endpoint
	GetTokenEndpoint               endpoint.Endpoint
}

// MakeServerEndpoints returns a struct containing all the endpoints for a PodcastManageService
func MakeServerEndpoints(svc PodcastManageService) Endpoints {
	return Endpoints{
		CreateUserEndpoint:             MakeCreateUserEndpoint(svc),
		GetUserEndpoint:                MakeGetUserEndpoint(svc),
		GetPodcastDetailsEndpoint:      MakeGetPodcastDetailsEndpoint(svc),
		SubscribeEndpoint:              MakeSubscribeEndpoint(svc),
		GetUserSubscriptionsEndpoint:   MakeGetUserSubscriptionsEndpoint(svc),
		GetSubscriptionDetailsEndpoint: MakeGetSubscriptionDetailsEndpoint(svc),
	}
}

// MakeGetSubscriptionDetailsEndpoint returns a GetSubscriptionDetailsEndpoint via the passed service
func MakeGetSubscriptionDetailsEndpoint(svc PodcastManageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getSubscriptionDetailsRequest)
		podcast, e := svc.GetSubscriptionDetails(ctx, req.EmailID, req.URL)
		if e != nil {
			return getSubscriptionDetailsResponse{podcast, err.Error()}, nil
		}
		return getSubscriptionDetailsResponse{podcast, ""}, nil
	}
}

// MakeGetUserSubscriptions returns an endpoint for getting user subscriptions via the passed service
func MakeGetUserSubscriptionsEndpoint(svc PodcastManageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUserSubscriptionsRequest)
		subscriptions, e := svc.GetUserSubscriptions(ctx, req.EmailID)
		if e != nil {
			return getUserSubscriptionsResponse{subscriptions, err.Error()}, nil
		}
		return getUserSubscriptionsResponse{subscriptions, ""}, nil
	}
}

// MakeCreateUserEndpoint returns a CreateUserEndpoint via the passed service
func MakeCreateUserEndpoint(svc PodcastManageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createUserRequest)
		e := svc.CreateUser(ctx, req.EmailID, req.Password)
		if e != nil {
			return createUserResponse{false, e.Error()}, nil
		}
		return createUserResponse{Status: true, Err: ""}, nil
	}
}

// MakeGetUserEndpoint returns a GetUserEndpoint via the passed service
func MakeGetUserEndpoint(svc PodcastManageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getUserRequest)
		user, e := svc.GetUser(ctx, req.EmailID)
		if e != nil {
			return getUserResponse{User: user, Err: e.Error()}, nil
		}
		return getUserResponse{user, ""}, nil
	}
}

// MakeGetPodcastDetailsEndpoint returns a GetPodcastDetailsEndpoint via the passed service
func MakeGetPodcastDetailsEndpoint(svc PodcastManageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getPodcastDetailsRequest)
		podcast, err := svc.GetPodcastDetails(ctx, req.URL)
		if err != nil {
			return getPodcastDetailsResponse{Podcast: podcast, Err: err.Error()}, nil
		}
		return getPodcastDetailsResponse{Podcast: podcast, Err: ""}, nil
	}
}

// MakeSubscribeEndpoint returns a SubscribeEndpoint via the passed service
func MakeSubscribeEndpoint(svc PodcastManageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(subscribeRequest)
		err = svc.Subscribe(ctx, req.EmailID, req.URL)
		if err != nil {
			return subscribeResponse{Status: false, Err: err.Error()}, nil
		}
		return subscribeResponse{Status: true, Err: ""}, nil
	}
}

type getSubscriptionDetailsRequest struct {
	EmailID string `json:"email_id"`
	URL     string `json:"url"`
}

type getSubscriptionDetailsResponse struct {
	Podcast podcastmg.Podcast `json:"podcast"`
	Err     string            `json:"err"`
}

type getUserSubscriptionsRequest struct {
	EmailID string `json:"email_id"`
}

type getUserSubscriptionsResponse struct {
	Subscriptions []podcastmg.Podcast `json:"subscriptions"`
	Error         string              `json:"error"`
}

type subscribeRequest struct {
	EmailID string `json:"email_id"`
	URL     string `json:"url"`
}

type subscribeResponse struct {
	Status bool   `json:"status"`
	Err    string `json:"err,omitempty"`
}

type getPodcastDetailsRequest struct {
	URL string `json:"url"`
}

type getPodcastDetailsResponse struct {
	Podcast podcastmg.Podcast
	Err     string `json:"err,omitempty"`
}

type getUserRequest struct {
	EmailID string `json:"email_id"`
}

type getUserResponse struct {
	User podcastmg.User `json:"user"`
	Err  string         `json:"err,omitempty"`
}

type createUserRequest struct {
	EmailID  string `json:"email_id"`
	Password string `json:"password"`
}

type createUserResponse struct {
	Status bool   `json:"status"`
	Err    string `json:"err,omitempty"`
}
