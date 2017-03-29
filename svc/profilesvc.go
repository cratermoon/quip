package svc

import (
	"context"
	"encoding/json"
	"errors"
	stdlog "log"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/expvar"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/cratermoon/quip/models"
	"github.com/cratermoon/quip/storage"
)

var profileCount metrics.Counter
var serviceHistogram metrics.Histogram

// ProfileService tells the time with many options
type ProfileService interface {
	GetProfile() models.Profile
	PostProfile(string, string)
}

type profileService struct {
	s storage.ProfileStorage
}

func (ps profileService) GetProfile(id string) (models.Profile, error) {
	return ps.s.Get(id)
}

func (ps profileService) PostProfile(name string, addr string) (models.Profile, error) {
	p, err := ps.s.Add(name, addr)
	return p, err
}

type profileRequest struct {
	ID string `json:"id"`
}

type profileResponse struct {
	Status string         `json:"status"`
	Person models.Profile `json:"person"`
}

type postProfileRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type postProfileResponse struct {
	Status string         `json:"status"`
	Person models.Profile `json:"person"`
}

func makeProfileEndpoint(ps profileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(profileRequest)
		p, err := ps.GetProfile(req.ID)
		if err != nil {
			return profileResponse{Status: err.Error(), Person: p}, nil
		}
		return profileResponse{Status: "ok", Person: p}, nil
	}
}

func makePostProfileEndpoint(ps profileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postProfileRequest)
		stdlog.Print("post profile")
		p, err := ps.PostProfile(req.Name, req.Address)
		if err != nil {
			return postProfileResponse{err.Error(), p}, nil
		}
		return postProfileResponse{"ok", p}, nil
	}
}

// NewProfileService initializes the Profile Service
func NewProfileService(r *mux.Router) {

	svc := profileService{storage.NewProfileStorage()}

	serviceHistogram = expvar.NewHistogram("profile_create_histogram", 50)
	profileCount = expvar.NewCounter("profile_count")

	metricsEndpoint := MakeMetricsMiddleware(serviceHistogram, profileCount)

	logger := kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout))
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger))
	loggingEndpoint := MakeLoggingMiddleware()

	getProfileHandler := httptransport.NewServer(
		makeProfileEndpoint(svc),
		decodeProfileRequest,
		encodeProfileResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	)

	postProfileHandler := httptransport.NewServer(
		loggingEndpoint(metricsEndpoint(makePostProfileEndpoint(svc))),
		decodePostProfileRequest,
		encodePostProfileResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	)

	r.Methods("GET").Path("/profile/{id}").Handler(getProfileHandler)
	r.Methods("POST").Path("/profile").Handler(postProfileHandler)
}

func decodeProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("No ID")
	}
	return profileRequest{id}, nil
}

func encodeProfileResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodePostProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req postProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}
	return req, nil
}

func encodePostProfileResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
