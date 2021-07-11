package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fangjjcs/bookings-app/pkg/config"
	"github.com/fangjjcs/bookings-app/pkg/driver"
	"github.com/fangjjcs/bookings-app/pkg/forms"
	"github.com/fangjjcs/bookings-app/pkg/helpers"
	"github.com/fangjjcs/bookings-app/pkg/models"
	"github.com/fangjjcs/bookings-app/pkg/render"
	"github.com/fangjjcs/bookings-app/pkg/repository"
	"github.com/fangjjcs/bookings-app/pkg/repository/dbrepo"
	"github.com/go-chi/chi"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL, a),
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

// PostSearchAvailability is for "book now" search
func (m *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w,err)
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w,err)
	}

	rooms, err := m.DB.SearchAvailibilityForAllRooms(startDate, endDate)
	if err != nil{
		helpers.ServerError(w,err)
		return 
	}

	for _ , r := range rooms {
		m.App.InfoLog.Println("ROOM:",r.ID, r.RoomName)
	}

	if len(rooms)==0{
		m.App.Session.Put(r.Context(),"error","No Availibility") 
		http.Redirect(w,r,"search-availability",http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservations{
		StartDate: startDate,
		EndDate: endDate,
	}
	// put these information into session in order to make a reservation in other page.
	m.App.Session.Put(r.Context(),"reservation",res)

	// send data to the template
	render.RenderTemplate(w, r,"choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})	
}

// SearchAvailability is the handler for the search-availability page
func (m *Repository) SearchAvailabilityByRoomID(w http.ResponseWriter, r *http.Request) {
	
	roomID, _ := strconv.Atoi(r.Form.Get("room_id")) 
	roomType := r.Form.Get("room_type")

	room := models.Room{
		ID: roomID,
		RoomName: roomType,
	}
	data := make(map[string]interface{})
	data["room"] = room
	m.App.Session.Put(r.Context(),"Room",room)
	// send data to the template
	render.RenderTemplate(w, r,"search-availability.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) PostSearchAvailabilityByRoomID(w http.ResponseWriter, r *http.Request) {
	
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	room, _ := m.App.Session.Get(r.Context(),"Room").(models.Room)
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w,err)
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w,err)
	}

	available, _ := m.DB.SearchAvailabilityByDatesAndRoomID(startDate, endDate,room.ID)
	log.Println(available)
	if !available{
		// log.Println("here")
		m.App.Session.Put(r.Context(),"error","Sorry, We don't have available room now.")
		log.Println(m.App.Session)
		if room.ID == 1{
			http.Redirect(w,r,"/generals",http.StatusSeeOther)
		}else{
			http.Redirect(w,r,"/majors",http.StatusSeeOther)
		}
		
		return
	}
	// log.Println("and here")
	res := models.Reservations{
		StartDate: startDate,
		EndDate: endDate,
		RoomID: room.ID,
		Room: room,
	}
	
	// put these information into session in order to make a reservation in other page.
	m.App.Session.Put(r.Context(),"reservation",res)
	http.Redirect(w,r,"/make-reservation",http.StatusSeeOther)	
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	
	// read id from url
	roomID, err := strconv.Atoi(chi.URLParam(r,"id")) //from url format "/choose-room/{id}"
	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	roomType := chi.URLParam(r,"room_name") //from url format "/choose-room/{id}"
	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	// Get information from session
	res, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservations)
	if !ok{
		helpers.ServerError(w,err)
		return
	}
	room := models.Room{
		RoomName: roomType,
	}
	
	// Add and Put information back into the session
	res.RoomID = roomID
	res.Room = room
	m.App.Session.Put(r.Context(),"reservation", res)
	http.Redirect(w,r,"/make-reservation",http.StatusSeeOther)

}

// MakeReservation is the handler for the MakeReservation page
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	// Get information from session
	res, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservations)
	if !ok{
		helpers.ServerError(w,errors.New("cannot get reservation from session"))
		return
	}
	
	// Reformat time from GO time format to "YYYY-MM-DD" and display on template instead of displaying time data in reservation
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	// send data to the template
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: stringMap,
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

	reservation, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservations)
	if !ok{
		helpers.ServerError(w,errors.New("cannot get reservation from session"))
		return
	}
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")
	
	// Check for validation
	form := forms.New(r.PostForm)
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
		return
	}

	// Booking by insert Reservation !
	newReservationID, err := m.DB.InsertReservations(reservation)
	if err != nil{
		helpers.ServerError(w,err)
	}
	// Adding a restriction after a successful booking 
	restriction := models.RoomRestrictions{
		StartDate: reservation.StartDate,
		EndDate: reservation.EndDate,
		RoomID: reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1, // 1 for reservation
	}
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil{
		helpers.ServerError(w,err)
	}


	// Build-up Confirm e-mail
	mailMsg := fmt.Sprintf(`
	 	<strong>Reservation Confirmation</strong><br>
		 <br>
		 Dear %s, <br>
		 This is a confirmation for your reservation from %s to %s.
	`,reservation.FirstName,reservation.StartDate.Format("2006-01-02"),reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To: reservation.Email,
		From: "server@booking.com",
		Subject: "Reservation Confirmation",
		Content: mailMsg,
	}
	// Put msg in the channel
	m.App.MailChan <- msg


	// transmit reservation data by session
	m.App.Session.Put(r.Context(),"reservation", reservation)
	http.Redirect(w,r,"reservation-summary", http.StatusSeeOther)
}

// ReservationSummary Get data from session and load into reservation-summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request){
	// get reservation data from a session
	reservation, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservations)
	// if there doesn't exist reservation obj in session
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from the session!")
		// For Alertion on the top
		m.App.Session.Put(r.Context(),"error", "Can not get data from the session!")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(),"reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed


	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
		StringMap: stringMap,
	})
}

// Contact is the handler for the about page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	// send data to the template
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{
	})
}

// Login page
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// Handles the Login
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	//whenever doing a login/logout => renew a token
	_ = m.App.Session.RenewToken(r.Context())

	err:= r.ParseForm()
	if err != nil{
		log.Println(err)
	}

	
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// validation
	form := forms.New(r.PostForm)
	form.Required("email","password")
	form.IsEmail("email",r)

	if !form.Valid(){
		// rerender the login page
		render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email,password)
	// Login Failed
	if err != nil{
		log.Println(err)
		m.App.Session.Put(r.Context(),"error","Invalid login credentials")
		http.Redirect(w, r, "/user/login",http.StatusSeeOther)
		return
	}
	// Login successfully
	m.App.Session.Put(r.Context(),"user_id", id)
	m.App.Session.Put(r.Context(),"flash","Login Successfully")
	http.Redirect(w,r,"/",http.StatusSeeOther)
}

// Log out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request){

	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w,r,"/",http.StatusSeeOther)
	m.App.Session.Put(r.Context(),"flash","Logout Successfully")


}

//
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w,r,"admin-dashboard.page.tmpl", &models.TemplateData{})
}

// Get All New Reservations in admin tool
func (m *Repository) AdminNewReservation(w http.ResponseWriter, r *http.Request){
	reservations, err := m.DB.AllNewReservations()
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.RenderTemplate(w,r,"admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminAllReservation(w http.ResponseWriter, r *http.Request){
	reservations, err := m.DB.AllReservations()
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.RenderTemplate(w,r,"admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// Show the reservation detail in the admin tool
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request){
	
	//Get reservation from db
	exploded := strings.Split(r.RequestURI,"/")
	
	id, err := strconv.Atoi(exploded[4])
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	
	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	reservation, err := m.DB.GetReservationByID(id)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.RenderTemplate(w,r,"admin-reservation-show.page.tmpl", &models.TemplateData{
		Data: data,
		StringMap: stringMap,
		Form: forms.New(nil),
	})
}

// update reservation in admin tool (POST)
func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request){

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w,err) //Put error message to helpers and print it
		return
	}

	//Get reservation from db
	exploded := strings.Split(r.RequestURI,"/")
		
	id, err := strconv.Atoi(exploded[4])
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	reservation, err := m.DB.GetReservationByID(id)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

    err = m.DB.UpdateReservation(reservation,id)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	//GO-TO AdminShowReservation handler, rerender and show flash
	m.App.Session.Put(r.Context(),"flash","Changes saved!")
	http.Redirect(w,r, fmt.Sprintf("/admin/reservations-%s",src),http.StatusSeeOther)
}

func (m *Repository) AdminReservationCalendar(w http.ResponseWriter, r *http.Request){
	render.RenderTemplate(w,r,"admin-reservations-calendar.page.tmpl", &models.TemplateData{})
}

// Marks a reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request){
	
	id, _ := strconv.Atoi(chi.URLParam(r,"id"))
	src := chi.URLParam(r,"src")
	// These code also work 
	// exploded := strings.Split(r.RequestURI,"/")
	// id, err := strconv.Atoi(exploded[4])
	// if err != nil{
	// 	helpers.ServerError(w,err)
	// 	return
	// }
	// src := exploded[3]
	
	_ = m.DB.UpdateProcessedForReservation(id,1)

	m.App.Session.Put(r.Context(),"flash","Reservation is marked as Processed!")
	http.Redirect(w,r, fmt.Sprintf("/admin/reservations-%s",src),http.StatusSeeOther)

}

// Delete a reservation
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request){
	
	id, _ := strconv.Atoi(chi.URLParam(r,"id"))
	src := chi.URLParam(r,"src")
	
	_ = m.DB.DeleteReservation(id)

	m.App.Session.Put(r.Context(),"flash","Reservation is deleted.")
	http.Redirect(w,r, fmt.Sprintf("/admin/reservations-%s",src),http.StatusSeeOther)

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