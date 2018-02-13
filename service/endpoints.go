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
	UnsubscribeEndpoint            endpoint.Endpoint
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
		UnsubscribeEndpoint:            MakeUnsubscribeEndpoint(svc),
		GetUserSubscriptionsEndpoint:   MakeGetUserSubscriptionsEndpoint(svc),
		GetSubscriptionDetailsEndpoint: MakeGetSubscriptionDetailsEndpoint(svc),
		GetTokenEndpoint:               MakeGetTokenEndpoint(svc),
	}
}

// MakeGetTokenEndpoint returns a GetTokenEndpoint via the passed service
func MakeGetTokenEndpoint(svc PodcastManageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getTokenRequest)
		tokenString, e := svc.GetToken(ctx, req.EmailID, req.Password)
		if e != nil {
			return getTokenResponse{tokenString, e.Error()}, e
		}
		return getTokenResponse{tokenString, ""}, nil
	}
}

// MakeGetSubscriptionDetailsEndpoint returns a GetSubscriptionDetailsEndpoint via the passed service
func MakeGetSubscriptionDetailsEndpoint(svc PodcastManageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getSubscriptionDetailsRequest)
		podcast, e := svc.GetSubscriptionDetails(ctx, req.EmailID, req.URL)
		if e != nil {
			return getSubscriptionDetailsResponse{podcast, e.Error()}, e
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
			return getUserSubscriptionsResponse{subscriptions, e.Error()}, e
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
			return createUserResponse{false, e.Error()}, e
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
			return getUserResponse{User: user, Err: e.Error()}, e
		}
		return getUserResponse{user, ""}, nil
	}
}

// MakeGetPodcastDetailsEndpoint returns a GetPodcastDetailsEndpoint via the passed service
func MakeGetPodcastDetailsEndpoint(svc PodcastManageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getPodcastDetailsRequest)
		podcast, e := svc.GetPodcastDetails(ctx, req.URL)
		if e != nil {
			return getPodcastDetailsResponse{Podcast: podcast, Err: err.Error()}, e
		}
		return getPodcastDetailsResponse{Podcast: podcast, Err: ""}, nil
	}
}

// MakeSubscribeEndpoint returns a SubscribeEndpoint via the passed service
func MakeSubscribeEndpoint(svc PodcastManageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(subscribeRequest)
		e := svc.Subscribe(ctx, req.EmailID, req.URL)
		if e != nil {
			return subscribeResponse{Status: false, Err: e.Error()}, e
		}
		return subscribeResponse{Status: true, Err: ""}, nil
	}
}

// MakeUnsubscribeEndpoint returns an UnsubscribeEndpoint via the passed service
func MakeUnsubscribeEndpoint(svc PodcastManageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(unsubscribeRequest)
		e := svc.Unsubscribe(ctx, req.EmailID, req.URL)
		if e != nil {
			return unsubscribeResponse{Status: false, Err: e.Error()}, e
		}
		return unsubscribeResponse{Status: true, Err: ""}, nil
	}
}

type getTokenRequest struct {
	EmailID  string `json:"email_id"`
	Password string `json:"password"`
}

type getTokenResponse struct {
	TokenString string `json:"token_string"`
	Err         string `json:"err,omitempty"`
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

type unsubscribeRequest struct {
	EmailID string `json:"email_id"`
	URL     string `json:"url"`
}

type unsubscribeResponse struct {
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
