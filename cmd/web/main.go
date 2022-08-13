package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/jkroller/bookings/pkg/config"
	"github.com/jkroller/bookings/pkg/handlers"
	"github.com/jkroller/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

var portNumber = ":8090"
var app config.AppConfig
var session *scs.SessionManager

// main is main application
func main() {

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // in production should be true

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Printf("Starting application on port %s\n", portNumber)
	// _ = is don't care error
	//_ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
