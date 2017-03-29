package svc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/expvar"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/cratermoon/quip/models"
	"github.com/cratermoon/quip/storage"
)

var profileCount metrics.Counter

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
		return profileResponse{p}, err
	}
}

func makePostProfileEndpoint(ps profileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(postProfileRequest)
		p, err := ps.PostProfile(req.Name, req.Address)
		if err != nil {
			return postProfileResponse{err.Error(), p}, err
		}
		profileCount.Add(1)
		return postProfileResponse{"ok", p}, nil
	}
}

// NewProfileService initializes the Profile Service
func NewProfileService(r *mux.Router) {

	svc := profileService{storage.NewProfileStorage()}

	getProfileHandler := httptransport.NewServer(
		makeProfileEndpoint(svc),
		decodeProfileRequest,
		encodeProfileResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	)

	postProfileHandler := httptransport.NewServer(
		makePostProfileEndpoint(svc),
		decodePostProfileRequest,
		encodePostProfileResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	)

	r.Methods("GET").Path("/profile/{id}").Handler(getProfileHandler)
	r.Methods("POST").Path("/profile").Handler(postProfileHandler)

	profileCount = expvar.NewCounter("profile_count")
}

func decodeProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	log.Println("Method", ctx.Value(httptransport.ContextKeyRequestMethod))
	log.Println("Path", ctx.Value(httptransport.ContextKeyRequestPath))
	return profileRequest{}, nil
}

func encodeProfileResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodePostProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return postProfileRequest{}, nil
}

func encodePostProfileResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
