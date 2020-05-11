package main

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"net"
	"net/http"
	"os"
	"os/signal"
	"test/api/server"
	"time"
)
var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{"user_id": 1})
	fmt.Printf("DEBUG: a sample jwt is %s\n", tokenString)
}

func NewRouter() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(render.SetContentType(render.ContentTypeJSON))
	//router.Mount("/", templates.Router())
	router.Use(jwtauth.Verifier(tokenAuth))
	router.Use(jwtauth.Authenticator)
	router.Mount("/api/v1/", server.Router())
	return router
}

func main() {

	server.DbStart()
	handler := NewRouter()
	srv := &http.Server{
		Handler: handler,
	}
	lr, _ := net.Listen("tcp", ":8080")
	go func() {
		fmt.Print("Server is running!\n\n")
		err := srv.Serve(lr)
		if err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

}
