package svc

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/cratermoon/quip/proto"
	"github.com/cratermoon/quip/uuid"
)

// UUIDService returns a randomly-generated UUID
type UUIDService interface {
	GetUUID() (string, error)
}

type uuidService struct{}

func (u uuidService) GetUUID() (string, error) {
	return uuid.NewUUID()
}

type uuidRequest struct{}

func makeUUIDEndpoint(us uuidService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		u, err := us.GetUUID()
		if err != nil {
			return proto.UUIDResponse{"err " + err.Error(), u}, nil
		}
		return proto.UUIDResponse{"ok", u}, nil
	}
}

// NewUUIDService initializes the UUID Service
func NewUUIDService(r *mux.Router) {

	svc := uuidService{}

	storeMiddleware := MakeStorageMiddleware()
	uuidhandler := httptransport.NewServer(
		storeMiddleware(makeUUIDEndpoint(svc)),
		decodeUUIDRequest,
		encodeUUIDResponse,
	)

	r.Methods("GET").Path("/uuid").Handler(uuidhandler)
}

func decodeUUIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return uuidRequest{}, nil
}

func encodeUUIDResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
