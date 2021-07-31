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

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["year"] = year
	stringMap["month"] = month

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

	m.App.Session.Put(r.Context(),"flash","Changes saved!")

	month := r.Form.Get("month")
	year := r.Form.Get("year")

	// GO-TO AdminShowReservation handler, rerender and show flash
	// if year != "" means we come from calendar page, so when we post form we need to go back to calendar page
	if year == ""{
		http.Redirect(w,r, fmt.Sprintf("/admin/reservations-%s",src),http.StatusSeeOther)
	}else{
		http.Redirect(w,r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s",year, month),http.StatusSeeOther)
	}
}



// AdminProcessReservation Marks a reservation as processed
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

	month := r.URL.Query().Get("m")
	year := r.URL.Query().Get("y")

	if year == ""{
		http.Redirect(w,r, fmt.Sprintf("/admin/reservations-%s",src),http.StatusSeeOther)
	}else{
		http.Redirect(w,r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s",year, month),http.StatusSeeOther)
	}
	

}

//AdminDeleteReservation Delete a reservation
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request){
	
	id, _ := strconv.Atoi(chi.URLParam(r,"id"))
	src := chi.URLParam(r,"src")
	
	_ = m.DB.DeleteReservation(id)

	month := r.URL.Query().Get("m")
	year := r.URL.Query().Get("y")


	m.App.Session.Put(r.Context(),"flash","Reservation is deleted.")

	if year == ""{
		http.Redirect(w,r, fmt.Sprintf("/admin/reservations-%s",src),http.StatusSeeOther)
	}else{
		http.Redirect(w,r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s",year, month),http.StatusSeeOther)
	}
}

//AdminReservationCalendar shows reservation on calendar
func (m *Repository) AdminReservationCalendar(w http.ResponseWriter, r *http.Request){

	// Assume that there is no year/month specified
	now := time.Now()
	
	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month),1,0,0,0,0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0,1,0) //next month
	last := now.AddDate(0,-1,0) //last month

	nextMonth := next.Format("01")
	nextYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	// how to get days in a month
	// get the first and the last day in month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth,1,0,0,0,0,currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0,1,-1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	data["rooms"] = rooms

	//get restrictions
	//GetRestrictionsForRoomByDate()
	for _, x := range rooms{
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d:= firstOfMonth;!d.After(lastOfMonth) ; d=d.AddDate(0,0,1){
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		//get restrictions for current room
		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID,firstOfMonth,lastOfMonth)
		if err != nil {
			helpers.ServerError(w,err)
			return
		}

		for _, y := range restrictions {
			if y.ReservationID > 0 {
				//it's a reservation
				for d := y.StartDate; !d.After(y.EndDate); d = d.AddDate(0,0,1){
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			}else{
				// it's a block
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID
			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap
		
		// store block_map n session
		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID),blockMap)
		
		//log.Println(blockMap)
		// 0 : available
		// 1 : Reserved
		// 2 : block
	}

	
	
	render.RenderTemplate(w,r,"admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		IntMap: intMap,
		Data: data,
	})

	m.App.Session.Remove(r.Context(),"flash")

}

//AdminPostReservationCalendar handles post of reservation calendar (change and save)
func (m *Repository) AdminPostReservationCalendar(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	//-- process blocks --//
	rooms, err := m.DB.AllRooms()
	if err != nil{
		helpers.ServerError(w,err)
		return
	}

	// handle delete blocks
	form := forms.New(r.PostForm)
	for _, x := range rooms {
		// Get the block map from the session. Loop through entire map, if we have an entry in the map
		// that does not exist in our posted data, and if the restriction id > 0, then it is a block we 
		// need to remove.

		// 1. get the block map from session
		curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		for name, value := range curMap{
			if val, ok := curMap[name]; ok { // ok will be false if the value not in the map
				// only pay attention to value > 0, and that are not in the form post
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name),r){
						log.Println("DELETE block", value)

						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	}

	// handle new blocks
	for name, _ := range r.PostForm { // name, value => ignore value
		//log.Println("Forms has name", name)
		
		// when the name of post-form has prefix of add_block means we add a block.
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name,"_")
			roomID, _ := strconv.Atoi(exploded[2])
			t, _ := time.Parse("2006-01-2", exploded[3])
			// insert a new block
			log.Println("INSERT a block for date", exploded[3])
			err := m.DB.InsertBlockForRoom(roomID,t)
			if err != nil {
				log.Println(err)
			}

		}

	}

	m.App.Session.Put(r.Context(),"flash","Changes Saved!")
	http.Redirect(w,r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year,month), http.StatusSeeOther)

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