package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fangjjcs/bookings-app/pkg/config"
	"github.com/fangjjcs/bookings-app/pkg/forms"
	"github.com/fangjjcs/bookings-app/pkg/helpers"
	"github.com/fangjjcs/bookings-app/pkg/models"
	"github.com/fangjjcs/bookings-app/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send data to the template
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Generals is the handler for the about page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	// send data to the template
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{
	})
}
// Majors is the handler for the about page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {

	// send data to the template
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{
	})
}
// SearchAvailability is the handler for the about page
func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {

	// send data to the template
	render.RenderTemplate(w, r,"search-availability.page.tmpl", &models.TemplateData{
	})
}
// Post : SearchAvailability
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("Search Availability from %s to %s",start,end)))
}

// JSON Response
type JsonResponse struct{
	STATUS    bool `json:"result"`
	MESSAGE   string `json:"msg"`
}
func (m *Repository) JsonSearchAvailability(w http.ResponseWriter, r *http.Request) {
	
	resp := JsonResponse{
		STATUS : true,
		MESSAGE : "Available!",
	}
	// JsonResponse -> []byte
	jsonResult, err := json.MarshalIndent(resp,"","    ")
	if err != nil {
		// log.Println(err)
		helpers.ServerError(w,err)
	}
	log.Println(string(jsonResult))
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)

}


// MakeReservation is the handler for the MakeReservation page
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation;
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	// send data to the template
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostMakeReservation is the handler for Post the reservation form
func (m *Repository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// log.Println(err)
		helpers.ServerError(w,err) //Put error message to helpers and print it
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName: r.Form.Get("last_name"),
		Email: r.Form.Get("email"),
		Phone: r.Form.Get("phone"),
	}

	// log.Println(r.PostForm)
	form := forms.New(r.PostForm)
	// form.Has("first_name", r)
	form.Required("first_name","last_name","email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email",r)

	if !form.Valid(){
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		// log.Println(form)
		return
	}
	// log.Println(r.Context())
	m.App.Session.Put(r.Context(),"reservation", reservation)
	http.Redirect(w,r,"reservation-summary", http.StatusSeeOther)
}

// Contact is the handler for the about page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	// send data to the template
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{
	})
}


// ReservationSummary Get data from session and load into reservation-summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request){
	reservation, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	// if there doesn't exist reservation obj in session
	if !ok {
		// log.Println("Can not get data from the session!")
		
		m.App.ErrorLog.Println("Can't get data from the session!")
		// For Alertion on the top
		m.App.Session.Put(r.Context(),"error", "Can not get data from the session!")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(),"reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}