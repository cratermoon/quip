package server

import (
	"expvar"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cratermoon/quip/svc"
)

// BuildServices creates the endpoints
func BuildServices() http.Handler {
	r := mux.NewRouter()
	svc.NewQuipService(r, "quip.crt")
	svc.NewTimeService(r)
	svc.NewUUIDService(r)

	r.Methods("GET").Path("/debug/vars").Handler(expvar.Handler())

	fmt.Println("All services ready")
	return r
}
