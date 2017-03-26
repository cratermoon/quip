package svc

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/cratermoon/quip/quipdb"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// QuipService provides a quip server
type QuipService interface {
	GetQuip() (string, error)
	CountQuips() (int64, error)
}

type quipService struct {
	repo quipdb.QuipRepo
}

func (q quipService) GetQuip() (string, error) {
	return q.repo.Quip()
}

func (q quipService) CountQuips() (int64, error) {
	return q.repo.Count()
}

type getquipRequest struct{}

type getquipResponse struct {
	Quip string `json:"quip"`
	Err  string `json:"err,omitempty"`
}

type countRequest struct{}

type countResponse struct {
	Count int64  `json:"i"`
	Err   string `json:"err,omitempty"`
}

func makeCountEndpoint(qs QuipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		c, err := qs.CountQuips()
		if err != nil {
			return countResponse{c, err.Error()}, err
		}
		return countResponse{c, ""}, nil
	}
}

func makeGetQuipEndpont(qs QuipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		q, err := qs.GetQuip()
		if err != nil {
			return getquipResponse{q, err.Error()}, err
		}
		return getquipResponse{q, ""}, nil
	}

}

// Setup initializes the QuipService
func Setup() {

	svc := quipService{quipdb.NewQuipRepo()}

	quiphandler := httptransport.NewServer(
		makeGetQuipEndpont(svc),
		decodeRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeRequest,
		encodeResponse,
	)

	http.Handle("/quip", quiphandler)
	http.Handle("/count", countHandler)
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return countRequest{}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
