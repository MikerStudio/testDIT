package server

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

func UsersRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/{userID}/", func(r chi.Router) {
		r.Use(UserMDLWR)
		r.Get("/", GetUserData)
		r.Put("/", UpdateUser)
		r.Delete("/", DeleteUser)
	})

	r.Get("/", GetUsersList)
	r.Post("/", CreateUser)
	return r
}

func UserMDLWR(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user *User
		var err error
		var userIdInt int
		userID := chi.URLParam(r, "userID")
		if userID != "" {
			userIdInt, err = strconv.Atoi(userID)
			user, err = dbGetUser(userIdInt)
		} else {
			render.Render(w, r, &DefaultResponse{HTTPStatusCode: 404, StatusText: "Resource not found."})
			return
		}
		if err != nil {
			render.Render(w, r, &DefaultResponse{HTTPStatusCode: 404, StatusText: "Resource not found."})
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
