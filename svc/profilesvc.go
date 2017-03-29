package svc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/cratermoon/quip/models"
	"github.com/cratermoon/quip/uuid"
)

// ProfileService tells the time with many options
type ProfileService interface {
	GetProfile() models.Profile
	PostProfile(string, string)
}

type profileService struct{}

func (q profileService) GetProfile() (models.Profile, error) {
	u, err := uuid.NewUUID()
	return models.Profile{ID: u, Name: "", Addresses: nil}, err
}

func (q profileService) PostProfile(models.Profile) {

}

type profileRequest struct{}

type profileResponse struct {
	Person models.Profile `json:"person"`
}

type postProfileRequest struct {
	Name    string
	Address string
}

type postProfileResponse struct {
	Status string `json:"status"`
}

func makeProfileEndpoint(ps profileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		p, err := ps.GetProfile()
		return profileResponse{p}, err
	}
}

func makePostProfileEndpoint(ps profileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		l := models.Location("1 Paradise Lane")
		u, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		ps.PostProfile(models.Profile{ID: u, Name: "Adam", Addresses: []models.Location{l}})
		return postProfileResponse{"ok"}, nil
	}
}

// NewProfileService initializes the Profile Service
func NewProfileService(r *mux.Router) {

	svc := profileService{}

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
