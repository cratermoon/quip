package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"net/http"

	"github.com/cratermoon/quip/quipdb"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/expvar"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var quipsServed metrics.Counter
var quipLatency metrics.Histogram

// QuipService provides a quip server
type QuipService interface {
	GetQuip() (string, error)
	CountQuips() (int64, error)
}

type quipService struct {
	repo quipdb.QuipRepo
}

func (q quipService) GetQuip() (string, error) {
	begin := time.Now()
	quip, err := q.repo.Quip()
	if err != nil {
		return "Ponder nothingness", err
	}
	quipsServed.Add(1)
	quipLatency.Observe(time.Since(begin).Seconds())
	return quip, nil
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
			return countResponse{c, fmt.Sprintf("The Wisdom Service is unavailable: %s", err)}, nil
		}
		return countResponse{c, ""}, nil
	}
}

func makeGetQuipEndpont(qs QuipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		q, err := qs.GetQuip()
		if err != nil {
			return getquipResponse{q, fmt.Sprintf("The Wisdom Service is unavailable: %s", err)}, nil
		}
		return getquipResponse{q, ""}, nil
	}

}

// NewQuipService initializes the QuipService
func NewQuipService(r *mux.Router) {

	q, err := quipdb.NewQuipRepo()

	if err != nil {
		return
	}

	svc := quipService{q}

	quipsServed = expvar.NewCounter("quips_served")
	quipLatency = expvar.NewHistogram("quip_quickness", 50)

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

	r.Methods("GET").Path("/quip").Handler(quiphandler)
	r.Methods("GET").Path("/count").Handler(countHandler)
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return countRequest{}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
