package server

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

func EventsRouter () http.Handler {
	r := chi.NewRouter()
	r.Route("/{eventID}/", func(r chi.Router) {
		r.Use(EventMDLWR)
		r.Get("/", GetEventData)
		r.Put("/", UpdateEvent)
		r.Delete("/", DeleteEvent)
	})

	r.Get("/", GetEventList)
	r.Post("/", CreateEvent)
	return r
}

func EventMDLWR(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var event *Event
		var err error
		var eventIdInt int
		eventID := chi.URLParam(r, "eventID")
		if eventID != "" {
			eventIdInt, err = strconv.Atoi(eventID)
			event, err = dbGetEvent(eventIdInt)
		} else {
			render.Render(w, r, &DefaultResponse{HTTPStatusCode: 404, StatusText: "Resource not found."})
			return
		}
		if err != nil {
			render.Render(w, r, &DefaultResponse{HTTPStatusCode: 404, StatusText: "Resource not found."})
			return
		}

		ctx := context.WithValue(r.Context(), "event", event)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
