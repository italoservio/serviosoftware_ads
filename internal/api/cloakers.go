package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/italoservio/serviosoftware_ads/internal/deps"
	"github.com/italoservio/serviosoftware_ads/pkg/jwt"
)

func RegisterCloakersRoutes(mux *mux.Router, c *deps.Container) {
	mux.Handle(
		"/r/{encodedId}",
		http.HandlerFunc(c.CloakersHttpAPI.RedirectCloakerHttpAPI.RedirectCloaker)).
		Methods("GET")

	protected := mux.PathPrefix("/cloakers").Subrouter()
	protected.Use(jwt.Middleware(c.Env))

	protected.
		Handle(
			"",
			http.HandlerFunc(c.CloakersHttpAPI.CreateCloakerHttpAPI.CreateCloaker)).
		Methods("POST")
	protected.
		Handle(
			"",
			http.HandlerFunc(c.CloakersHttpAPI.ListCloakersHttpAPI.ListCloakers)).
		Methods("GET")
	protected.
		Handle(
			"/{cloakerId}",
			http.HandlerFunc(c.CloakersHttpAPI.UpdateCloakerHttpAPI.UpdateCloakerByID)).
		Methods("PUT", "PATCH")
	protected.
		Handle(
			"/{cloakerId}",
			http.HandlerFunc(c.CloakersHttpAPI.GetCloakerByIDHttpAPI.GetCloakerByID)).
		Methods("GET")
	protected.
		Handle(
			"/{cloakerId}",
			http.HandlerFunc(c.CloakersHttpAPI.DeleteCloakerHttpAPI.DeleteCloakerByID)).
		Methods("DELETE")
}
