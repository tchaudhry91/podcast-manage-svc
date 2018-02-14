package service

import (
	"context"
	"encoding/json"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	// ErrJSONUnmarshall is an error when the JSON parsing fails on the request
	ErrJSONUnmarshall = errors.New("Failed to parse incoming JSON")
)

// MakeHTTPHandler returns a router for the podcast-manager-service
func MakeHTTPHandler(svc PodcastManageService, signingString string, logger log.Logger) http.Handler {
	router := mux.NewRouter()
	endpoints := MakeServerEndpoints(svc)
	serverOptions := []kithttp.ServerOption{
		kithttp.ServerBefore(kitjwt.HTTPToContext()),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerErrorLogger(logger),
	}

	kf := func(token *jwt.Token) (interface{}, error) {
		return []byte(signingString), nil
	}
	claimsFetcher := func() jwt.Claims {
		return &TokenClaims{}
	}
	authMiddleware := kitjwt.NewParser(kf, jwt.SigningMethodHS256, claimsFetcher)

	router.Methods("POST").Path("/register").Handler(kithttp.NewServer(
		endpoints.CreateUserEndpoint,
		decodeCreateUserRequest,
		encodeGenericResponse,
		serverOptions...,
	))

	getUserEndpoint := endpoints.GetUserEndpoint
	getUserEndpoint = authMiddleware(getUserEndpoint)
	router.Methods("POST").Path("/user").Handler(kithttp.NewServer(
		getUserEndpoint,
		decodeGetUserRequest,
		encodeGenericResponse,
		serverOptions...,
	))
	router.Methods("GET").Path("/user/{user}").Handler(kithttp.NewServer(
		getUserEndpoint,
		decodeGetUserRequestAlternate,
		encodeGenericResponse,
		serverOptions...,
	))

	router.Methods("POST").Path("/podcast").Handler(kithttp.NewServer(
		endpoints.GetPodcastDetailsEndpoint,
		decodeGetPodcastDetailsRequest,
		encodeGenericResponse,
		serverOptions...,
	))

	subscribeEndpoint := endpoints.SubscribeEndpoint
	subscribeEndpoint = authMiddleware(subscribeEndpoint)
	router.Methods("POST").Path("/subscribe").Handler(kithttp.NewServer(
		subscribeEndpoint,
		decodeSubscribeRequest,
		encodeGenericResponse,
		serverOptions...,
	))

	unsubscribeEndpoint := endpoints.UnsubscribeEndpoint
	unsubscribeEndpoint = authMiddleware(unsubscribeEndpoint)
	router.Methods("POST").Path("/unsubscribe").Handler(kithttp.NewServer(
		unsubscribeEndpoint,
		decodeSubscribeRequest,
		encodeGenericResponse,
		serverOptions...,
	))

	subscriptionsEndpoint := endpoints.GetUserSubscriptionsEndpoint
	subscriptionsEndpoint = authMiddleware(subscriptionsEndpoint)
	router.Methods("POST").Path("/subscriptions").Handler(kithttp.NewServer(
		subscriptionsEndpoint,
		decodeGetUserSubscriptionsRequest,
		encodeGenericResponse,
		serverOptions...,
	))

	subscriptionEndpoint := endpoints.GetSubscriptionDetailsEndpoint
	subscriptionEndpoint = authMiddleware(subscriptionEndpoint)
	router.Methods("POST").Path("/subscription").Handler(kithttp.NewServer(
		subscriptionEndpoint,
		decodeGetSubscriptionDetailsRequest,
		encodeGenericResponse,
		serverOptions...,
	))

	router.Methods("POST").Path("/login").Handler(kithttp.NewServer(
		endpoints.GetTokenEndpoint,
		decodeGetTokenRequest,
		encodeGenericResponse,
		serverOptions...,
	))
	return router
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.WriteHeader(codeFrom(err))
	e := json.NewEncoder(w).Encode(map[string]interface{}{
		"err": err.Error(),
	})
	if e != nil {
		panic("Error encoding error")
	}
}

func decodeGetTokenRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var tokenReq getTokenRequest
	if err := json.NewDecoder(req.Body).Decode(&tokenReq); err != nil {
		return nil, ErrJSONUnmarshall
	}
	return tokenReq, nil
}

func decodeGetSubscriptionDetailsRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var subReq getSubscriptionDetailsRequest
	if err := json.NewDecoder(req.Body).Decode(&subReq); err != nil {
		return nil, ErrJSONUnmarshall
	}
	return subReq, nil
}

func decodeGetUserSubscriptionsRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var subReq getUserSubscriptionsRequest
	if err := json.NewDecoder(req.Body).Decode(&subReq); err != nil {
		return nil, ErrJSONUnmarshall
	}
	return subReq, nil
}

func decodeSubscribeRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var subReq subscribeRequest
	if err := json.NewDecoder(req.Body).Decode(&subReq); err != nil {
		return nil, ErrJSONUnmarshall
	}
	return subReq, nil
}

func decodeGetPodcastDetailsRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var podReq getPodcastDetailsRequest
	if err := json.NewDecoder(req.Body).Decode(&podReq); err != nil {
		return nil, ErrJSONUnmarshall
	}
	return podReq, nil
}

func decodeGetUserRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var guReq getUserRequest
	if err := json.NewDecoder(req.Body).Decode(&guReq); err != nil {
		return nil, ErrJSONUnmarshall
	}
	return guReq, nil
}

func decodeGetUserRequestAlternate(ctx context.Context, req *http.Request) (request interface{}, err error) {
	vars := mux.Vars(req)
	guReq := getUserRequest{
		EmailID: vars["user"],
	}
	return guReq, nil
}

func decodeCreateUserRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var userRequest createUserRequest
	if err := json.NewDecoder(req.Body).Decode(&userRequest); err != nil {
		return nil, ErrJSONUnmarshall
	}
	return userRequest, nil
}

func encodeGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func codeFrom(err error) int {
	switch err {
	case ErrJSONUnmarshall:
		return http.StatusBadRequest
	case ErrUserFetch:
		return http.StatusBadRequest
	case ErrInvalidPassword:
		return http.StatusUnauthorized
	case kitjwt.ErrTokenContextMissing:
		return http.StatusUnauthorized
	case kitjwt.ErrTokenInvalid:
		return http.StatusUnauthorized
	case kitjwt.ErrTokenExpired:
		return http.StatusUnauthorized
	case kitjwt.ErrTokenMalformed:
		return http.StatusBadRequest
	case kitjwt.ErrTokenNotActive:
		return http.StatusUnauthorized
	case ErrInvalidClaim:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
