package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/ikabir21/bookings/pkg/config"
	"github.com/ikabir21/bookings/pkg/handlers"
)

// METHOD - 1: using pat

// func routes(app *config.AppConfig) http.Handler {
// 	mux := pat.New()

// 	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
// 	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

// 	return mux
// }

// METHOD -2: using chi

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	// mux.Use(MyMiddleware)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux;
}