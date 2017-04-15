package svc

import (
	"context"
	"encoding/json"
	"time"

	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// TimeService tells the time with many options
type TimeService interface {
	CurrentTime() string
}

type timeService struct{}

func (q timeService) CurrentTime() string {
	return time.Now().UTC().Format(time.RFC1123Z)
}

type timeRequest struct{}

type timeResponse struct {
	Time string `json:"time"`
}

func makeTimeEndpoint(ts timeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return timeResponse{ts.CurrentTime()}, nil
	}
}

// NewTimeService initializes the Time Service
func NewTimeService(r *mux.Router) {

	svc := timeService{}

	timehandler := httptransport.NewServer(
		makeTimeEndpoint(svc),
		decodeTimeRequest,
		encodeTimeResponse,
	)

	r.Methods("GET").Path("/time").Handler(timehandler)
}

func decodeTimeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return timeRequest{}, nil
}

func encodeTimeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
