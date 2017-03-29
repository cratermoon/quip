package svc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/cratermoon/quip/models"
)

// ProfileService tells the time with many options
type ProfileService interface {
	GetProfile() models.Profile
	PostProfile(string, string)
}

type profileService struct{}

func (q profileService) GetProfile() models.Profile {
	return models.Profile{ID: "", Name: "", Addresses: nil}
}

type profileRequest struct{}

type profileResponse struct {
	Person models.Profile
}

func makeProfileEndpoint(ps profileService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		p := ps.GetProfile()
		return profileResponse{p}, nil
	}
}

// NewProfileService initializes the Profile Service
func NewProfileService() http.Handler {

	r := mux.NewRouter()
	svc := profileService{}

	profileHandler := httptransport.NewServer(
		makeProfileEndpoint(svc),
		decodeProfileRequest,
		encodeProfileResponse,
		httptransport.ServerBefore(httptransport.PopulateRequestContext),
	)

	r.Methods("GET").Path("/profile").Handler(profileHandler)
	return r
	//http.Handle("/profile/{id}", profileHandler)
}

func decodeProfileRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	log.Println("Method", ctx.Value(httptransport.ContextKeyRequestMethod))
	log.Println("Path", ctx.Value(httptransport.ContextKeyRequestPath))
	return profileRequest{}, nil
}

func encodeProfileResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
