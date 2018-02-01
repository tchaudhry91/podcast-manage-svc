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
	ErrJSONUnmarshall = errors.New("Failed to parse incoming JSON")
)

func MakeHTTPHandler(svc PodcastManageService) http.Handler {
	router := mux.NewRouter()
	endpoints := MakeServerEndpoints(svc)

	router.Methods("POST").Path("/register").Handler(kithttp.NewServer(
		endpoints.CreateUserEndpoint,
		decodeCreateUserRequest,
		encodeCreateUserResponse,
	))
	return router
}

func decodeCreateUserRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var userRequest createUserRequest
	if err := json.NewDecoder(req.Body).Decode(&userRequest); err != nil {
		return nil, ErrJSONUnmarshall
	}
	return userRequest, nil
}

func encodeCreateUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
