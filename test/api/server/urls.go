package server

import (
	"github.com/go-chi/chi"
	"net/http"
)


func Router() http.Handler {
	router := chi.NewRouter()
	router.Mount("/events/", EventsRouter())
	router.Mount("/companies/", CompaniesRouter())
	router.Mount("/users/", UsersRouter())
	return router
}