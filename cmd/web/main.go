package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/fangjjcs/bookings-app/pkg/config"
	"github.com/fangjjcs/bookings-app/pkg/driver"
	"github.com/fangjjcs/bookings-app/pkg/handlers"
	"github.com/fangjjcs/bookings-app/pkg/helpers"
	"github.com/fangjjcs/bookings-app/pkg/models"
	"github.com/fangjjcs/bookings-app/pkg/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {

	// test run()
	db, err := run()
	if err != nil{
		log.Fatal(err)
	}
	defer db.SQL.Close() 

	fmt.Printf(fmt.Sprintf("Staring application on port %s\n", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {

	// what am I going to put in the session
	gob.Register(models.Reservations{})
	// gob.Register(models.Reservations{})
	gob.Register(models.User{})
	gob.Register(models.Restrictions{})
	gob.Register(models.Room{})

	// change this to true when in production
	app.InProduction = false

	// log
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// fmt.Println(app.InfoLog)
	// fmt.Println(app.ErrorLog)

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=Bookings user=fang password=")
	if err != nil{
		log.Fatal("Cannot connect to database")
	}
	


	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false // define whenever you allow to use cache or not

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	return db, nil
}