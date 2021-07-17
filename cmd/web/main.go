package main

import (
	"encoding/gob"
	"flag"
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

const portNumber = ":8088"

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

	defer close(app.MailChan)
	listenForMail()
	fmt.Println("Starting mail listener...")


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
	gob.Register(models.User{})
	gob.Register(models.Restrictions{})
	gob.Register(models.Room{})
	gob.Register(map[string]int{})

	// read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser ==""{
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	
	// log
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	//"host=localhost port=5432 dbname=Bookings user=fang password="
	db, err := driver.ConnectSQL(connectionString)
	if err != nil{
		log.Fatal("Cannot connect to database")
	}
	


	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.InProduction =  *inProduction //true // change this to true when in production
	app.UseCache = *useCache // define whenever you allow to use cache or not

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	return db, nil
}