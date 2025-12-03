package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftware_ads/pkg/exception"
)

type Health struct {
	Server string `json:"server"`
}

func RegisterInfraRoutes(mux *mux.Router) {
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		res, _ := json.Marshal(Health{Server: "ok"})

		w.Write(res)
	})).Methods("GET")
}

func MethodNotAllowed(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	res, _ := json.Marshal(exception.AppException{
		Messages:   []string{"method not allowed"},
		StatusText: http.StatusText(http.StatusMethodNotAllowed),
		StatusCode: http.StatusMethodNotAllowed,
	})

	w.Write(res)
}
