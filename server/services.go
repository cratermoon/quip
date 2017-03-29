package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cratermoon/quip/svc"
)

func BuildServices() http.Handler {
	r := mux.NewRouter()
	svc.NewProfileService(r)
	svc.NewQuipService(r)
	svc.NewTimeService(r)

	fmt.Println("All services ready")
	return r
}
