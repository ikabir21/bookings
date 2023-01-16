package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("my middleware called...")
		next.ServeHTTP(w, r)
	})
}

// add CSRF protection to all POST request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// load and save the session on every request
func SessionLoad(next http.Handler) http.Handler{
	return session.LoadAndSave(next)
}
