package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router returns the keystore api's router
func Router() *mux.Router {
	r := mux.NewRouter()
	r.Methods(http.MethodPost).Path("/key").HandlerFunc(writeKey)
	r.Methods(http.MethodGet).Path("/key/{key_id}").HandlerFunc(readKey)
	r.Methods(http.MethodGet).Path("/keys").HandlerFunc(readKeys)
	return r
}
