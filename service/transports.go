package service

import (
	"context"
	"encoding/json"
	"errors"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	ErrJSONUnmarshall = errors.New("Failed to parse incoming JSON") // Indicates an unparsable JSON input
)

func MakeHTTPHandler(svc PodcastManageService) http.Handler {
	router := mux.NewRouter()
	endpoints := MakeServerEndpoints(svc)

	router.Methods("POST").Path("/register").Handler(kithttp.NewServer(
		endpoints.CreateUserEndpoint,
		decodeCreateUserRequest,
		encodeGenericResponse,
	))

	router.Methods("POST").Path("/user").Handler(kithttp.NewServer(
		endpoints.GetUserEndpoint,
		decodeGetUserRequest,
		encodeGenericResponse,
	))

	router.Methods("POST").Path("/podcast").Handler(kithttp.NewServer(
		endpoints.GetPodcastDetailsEndpoint,
		decodeGetPodcastDetailsRequest,
		encodeGenericResponse,
	))

	router.Methods("POST").Path("/subscribe").Handler(kithttp.NewServer(
		endpoints.SubscribeEndpoint,
		decodeSubscribeRequest,
		encodeGenericResponse,
	))
	return router
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
