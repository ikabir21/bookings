package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ikabir21/bookings/internal/config"
	"github.com/ikabir21/bookings/internal/handlers"
	"github.com/ikabir21/bookings/internal/render"
)

const PORT = ":8081"
var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	
	app.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc

	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Println("Server is running on http://localhost" + PORT)

	server := &http.Server{
		Addr:    PORT,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	log.Fatal(err)
}
