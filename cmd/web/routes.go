package main

import (
	"net/http"

	"github.com/fangjjcs/bookings-app/pkg/config"
	"github.com/fangjjcs/bookings-app/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals", handlers.Repo.Generals)
	mux.Get("/majors", handlers.Repo.Majors)
	
	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	mux.Get("/search-availability-json", handlers.Repo.JsonSearchAvailability)
	mux.Post("/search-availability-id",handlers.Repo.SearchAvailabilityByRoomID)
	mux.Post("/search-availability-by-id",handlers.Repo.PostSearchAvailabilityByRoomID)
	mux.Get("/choose-room/{id}-{room_name}",handlers.Repo.ChooseRoom)

	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Post("/make-reservation", handlers.Repo.PostMakeReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/user/login", handlers.Repo.Login)
	mux.Post("/user/login", handlers.Repo.PostLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)
	mux.Route("/admin", func(mux chi.Router){
		// mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard) 
		mux.Get("/reservations-new", handlers.Repo.AdminNewReservation)
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservation)
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationCalendar)

	})


	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
