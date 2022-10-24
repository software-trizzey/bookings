package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/IM-Deane/bookings/pkg/config"
	"github.com/IM-Deane/bookings/pkg/handlers"
	"github.com/IM-Deane/bookings/pkg/render"
)


const port = ":8080"
var app config.AppConfig
var session *scs.SessionManager

//  main is the app entry point
func main() {

	// FIXME: set to "true" in prod.
	app.InProduction = false

	// setup new session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	// persist session on browser close
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tempCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tempCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// pass app config
	render.NewTemplates(&app)

	fmt.Printf(fmt.Sprintf("Starting application on port %s", port))

	server := &http.Server{
		Addr: port,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	log.Fatal(err)
}

