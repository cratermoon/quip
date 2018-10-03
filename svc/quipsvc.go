package svc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/expvar"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/cratermoon/quip/proto"
	"github.com/cratermoon/quip/quipdb"
	"github.com/cratermoon/quip/signing"
	"github.com/cratermoon/quip/storage"
)

var quipsServed metrics.Counter
var quipLatency metrics.Histogram

// maxQuipLength is the longest quip allowed, leaving space for hashtag
const maxQuipLength = 274

// QuipService provides a quip server
type QuipService interface {
	Get() (string, error)
	Count() (int64, error)
	List() ([]string, error)
	Add(proto.AddQuipRequest) (string, error)
}

type quipService struct {
	repo quipdb.QuipRepo
	ver  signing.Verifier
}

type getRequest struct{}

type countRequest struct{}

type countResponse struct {
	Count int64  `json:"i"`
	Err   string `json:"err,omitempty"`
}

type listRequest struct{}

func (q quipService) Get() (string, error) {
	begin := time.Now()
	quip, err := q.repo.Quip()
	if err != nil {
		return "Ponder nothingness", err
	}
	quipsServed.Add(1)
	quipLatency.Observe(time.Since(begin).Seconds())
	return quip.Quip(), nil
}

func (q quipService) Count() (int64, error) {
	return q.repo.Count()
}

func (q quipService) List() ([]string, error) {
	return q.repo.List()
}

func (q quipService) Add(req proto.AddQuipRequest) (string, error) {
	quip := req.Quip
	if len(quip) > maxQuipLength {
		return "err", fmt.Errorf(
			"maximum quip length (%d) exceeded, got %d",
			maxQuipLength, len(quip))
	}
	uuid := req.UUID
	if uuid == "" {
		log.Printf("Empty UUID\n")
		return quip, fmt.Errorf("empty UUID")
	}
	log.Printf("Checking signature on %s:%s\n", quip, uuid)
	v := strings.Join([]string{quip, uuid}, ":")
	err := q.ver.Verify(v, req.Signature)
	if err != nil {
		log.Printf("Signature error (%q) %s\n", quip, err.Error())
		return quip, fmt.Errorf("signature error")
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
			return proto.QuipResponse{q, fmt.Sprintf("The Wisdom Service is unavailable: %s", err)}, nil
		}
		return proto.QuipResponse{q, ""}, nil
	}

}

func makeListEndpoint(qs QuipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		q, err := qs.List()
		if err != nil {
			return proto.ListQuipResponse{q, fmt.Sprintf("The Wisdom Service is unavailable: %s", err)}, nil
		}
		return proto.ListQuipResponse{q, ""}, nil
	}
}

func makeAddEndpoint(qs QuipService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(proto.AddQuipRequest)
		if !ok {
			return proto.AddQuipResponse{"", fmt.Sprintf("Experience tranquility")}, errors.New("type assertion failed")
		}
		name, err := qs.Add(req)
		if err != nil {
			return proto.AddQuipResponse{name, err.Error()}, nil
		}
		log.Printf("New quip added: %q\n", req.Quip)
		return proto.AddQuipResponse{name, ""}, nil
	}
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return countRequest{}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req proto.AddQuipRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeAddResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp, ok := response.(proto.AddQuipResponse)
	if !ok {
		return errors.New("Response type assertion failed")
	}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(resp)
}

// NewQuipService initializes the QuipService
func NewQuipService(r *mux.Router, keyFName string) {

	q, err := quipdb.NewQuipRepo()

	if err != nil {
		return
	}

	kit, err := storage.NewKit()
	if err != nil {
		log.Println("Error starting quip service", err)
		return
	}
	crt, err := kit.FileObject(keyFName)

	if err != nil {
		log.Println("Error starting quip service", err)
		return
	}

	v := signing.Verifier{Cert: crt}
	svc := quipService{q, v}

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

	lookupMiddleware := MakeLookupMiddleware()

	addhandler := httptransport.NewServer(
		lookupMiddleware(makeAddEndpoint(svc)),
		decodeAddRequest,
		encodeAddResponse,
	)

	r.Methods("GET").Path("/quip").Handler(gethandler)
	r.Methods("GET").Path("/list").Handler(listhandler)
	r.Methods("GET").Path("/count").Handler(countHandler)
	r.Methods("POST").Path("/quip").Handler(addhandler)
}
