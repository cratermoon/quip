package svc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/expvar"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/cratermoon/quip/quipdb"
	"github.com/cratermoon/quip/signing"
)

var quipsServed metrics.Counter
var quipLatency metrics.Histogram

const MAX_QUIP_LENGTH = 134

// QuipService provides a quip server
type QuipService interface {
	Get() (string, error)
	Count() (int64, error)
	List() ([]string, error)
	Add(string) (string, error)
}

type quipService struct {
	repo quipdb.QuipRepo
}

type getRequest struct{}

type getResponse struct {
	Quip string `json:"quip"`
	Err  string `json:"err,omitempty"`
}

type countRequest struct{}

type countResponse struct {
	Count int64  `json:"i"`
	Err   string `json:"err,omitempty"`
}

type listRequest struct{}

type listResponse struct {
	Quips []string `json:"quips"`
	Err   string   `json:"err,omitempty"`
}

type addRequest struct {
	Quip      string `json:"quip"`
	Signature string `json:"sig"`
}

type addResponse struct {
	Name string `json:"name"`
	Err  string `json:"err,omitempty"`
}

func (q quipService) Get() (string, error) {
	begin := time.Now()
	quip, err := q.repo.Quip()
	if err != nil {
		return "Ponder nothingness", err
	}
	quipsServed.Add(1)
	quipLatency.Observe(time.Since(begin).Seconds())
	return quip, nil
}

func (q quipService) Count() (int64, error) {
	return q.repo.Count()
}

func (q quipService) List() ([]string, error) {
	return q.repo.List()
}

func (q quipService) Add(quip string) (string, error) {
	if len(quip) > MAX_QUIP_LENGTH {
		return "too long", fmt.Errorf(
			"Maximum quip length (%d) exceeded, got %d",
			MAX_QUIP_LENGTH, len(quip))
	}
	return q.repo.Add(quip)
}

func makeCountEndpoint(qs QuipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		c, err := qs.Count()
		if err != nil {
			return countResponse{c, fmt.Sprintf("The Wisdom Service is unavailable: %s", err)}, nil
		}
		return countResponse{c, ""}, nil
	}
}

func makeGetEndpont(qs QuipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		q, err := qs.Get()
		if err != nil {
			return getResponse{q, fmt.Sprintf("The Wisdom Service is unavailable: %s", err)}, nil
		}
		return getResponse{q, ""}, nil
	}

}

func makeListEndpoint(qs QuipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		q, err := qs.List()
		if err != nil {
			return listResponse{q, fmt.Sprintf("The Wisdom Service is unavailable: %s", err)}, nil
		}
		return listResponse{q, ""}, nil
	}
}

func makeAddEndpoint(qs QuipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(addRequest)
		if ok {
			err := signing.Verify(req.Quip, req.Signature)
			if err != nil {
				return addResponse{"err", err.Error()}, nil
			}
			name, err := qs.Add(req.Quip)
			if err != nil {
				return addResponse{name, err.Error()}, nil
			}
			return addResponse{name, ""}, nil
		}
		return addResponse{"", fmt.Sprintf("Experience tranquility")}, errors.New("type assertion failed")
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return countRequest{}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req addRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeAddResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp, ok := response.(addResponse)
	if !ok {
		return errors.New("Response type assertion failed")
	}
	return json.NewEncoder(w).Encode(resp)
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

	gethandler := httptransport.NewServer(
		makeGetEndpont(svc),
		decodeRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeRequest,
		encodeResponse,
	)

	listhandler := httptransport.NewServer(
		makeListEndpoint(svc),
		decodeRequest,
		encodeResponse,
	)

	addhandler := httptransport.NewServer(
		makeAddEndpoint(svc),
		decodeAddRequest,
		encodeAddResponse,
	)

	r.Methods("GET").Path("/quip").Handler(gethandler)
	r.Methods("GET").Path("/list").Handler(listhandler)
	r.Methods("GET").Path("/count").Handler(countHandler)
	r.Methods("POST").Path("/quip").Handler(addhandler)
}
