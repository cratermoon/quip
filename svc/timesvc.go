package svc

import (
	"context"
	"encoding/json"
	"time"

	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
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

// TimeSetup initializes the Time Service
func TimeSetup() {

	svc := timeService{}

	timehandler := httptransport.NewServer(
		makeTimeEndpoint(svc),
		decodeTimeRequest,
		encodeTimeResponse,
	)

	http.Handle("/time", timehandler)
}

func decodeTimeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return timeRequest{}, nil
}

func encodeTimeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
