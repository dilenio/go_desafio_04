package main

import (
	"net/http"

	"github.com/dilenio/desafio04/cmd/configs"
	"github.com/dilenio/desafio04/limiter"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	rateLimiter := limiter.NewRateLimiter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(rateLimiter.LimitHandler)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	http.ListenAndServe("127.0.0.1:"+configs.WebServerPort, router)
}
