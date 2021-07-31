package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/fangjjcs/bookings-app/pkg/models"
	"golang.org/x/crypto/bcrypt"
)
func (m *postgresDBRepo) AllUsers() bool{
	return true
}

// InsertReservations inserts reservation
func (m *postgresDBRepo) InsertReservations(res models.Reservations) (int, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int
	stmt := `insert into reservations 
	(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
	 values($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	
    err := m.DB.QueryRowContext(ctx, stmt,
	res.FirstName,
	res.LastName,
	res.Email,
	res.Phone,
	res.StartDate,
	res.EndDate,
	res.RoomID,
	time.Now(),
	time.Now(),
	).Scan(&newID)

	if err != nil{
		return 0, err
	}
	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into database
func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestrictions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()


	stmt := `insert into room_restrictions 
	(start_date, end_date, room_id, reservation_id,created_at, updated_at, restriction_id)
	 values($1, $2, $3, $4, $5, $6, $7)`
	
	 _, err := m.DB.ExecContext(ctx, stmt,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		time.Now(),
		time.Now(),
		res.RestrictionID,
		)
	if err!=nil{
		return err
	}

	return nil
}


// SearchAvailabilityByDatesAndRoomID search availability by room id
func (m *postgresDBRepo) SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool,error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// log.Println(start,end, roomID)
	var count int
	stmt := `select count(id) from room_restrictions where
	        $1 < end_date and $2 > start_date and room_id = $3;`
	
	err := m.DB.QueryRowContext(ctx, stmt , start, end, roomID).Scan(&count)
	if err!=nil{
		return false, err
	}

	// log.Println(count)
	// assumes there exists only 1 room
	if count == 0 {
		return true, nil
	}
	return false, nil
}

// SearchAvailibilityForAllRooms handles room availability by dates
func (m *postgresDBRepo) SearchAvailibilityForAllRooms(start, end time.Time) ([]models.Room, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query :=`
	select
		r.id, r.room_name
	from
		rooms r 
	where r.id not in 
		(select rr.room_id from room_restrictions rr where rr.start_date <$2 and rr.end_date>$1)
		`
	rows, err := m.DB.QueryContext(ctx,query,start,end)
	if err != nil{
		return nil, err
	}
	var room models.Room
	for rows.Next(){
		err := rows.Scan(&room.ID,&room.RoomName)
		if err!= nil{
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

// GetUserByID gets user information by user ID
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level created_at, updated_at,
				from users where id=$1`
	
	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.AccessLevel, &u.CreatedAt, &u.UpdatedAt)
	if err != nil{
		return u, err
	}
	return u, nil
}

// UpdateUser updates user information
func (m *postgresDBRepo) UpdateUser(u models.User) (error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name=$1, last_name=$2, email=$3, access_level=$4, updated_at=$5`
	
	_, err := m.DB.ExecContext(ctx, query,u.FirstName,u.LastName,u.Email,u.AccessLevel,time.Now())
	if err != nil{
		return err
	}

	return nil
}

// Authenticate validates login information
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	// to see if this user exist in db or not
	query := `select id, password from users where email=$1`
	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil{
		return id, hashedPassword, err
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword{
		return 0, "", errors.New("Incorrect Password")
	}else if err != nil{
		return 0, "", err
	}

	return id, hashedPassword, nil

}

// AllReservations returns a slice of all reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservations, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservations

	query := ` select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
				r.end_date, r.room_id, r.created_at, r.updated_at, rm.id, rm.room_name
				from reservations r
				left join rooms rm on (r.room_id = rm.id)
				order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	
	for rows.Next(){
		var item models.Reservations
		err := rows.Scan(&item.ID,&item.FirstName,&item.LastName,&item.Email,&item.Phone,&item.StartDate,
			&item.EndDate,&item.RoomID,&item.CreatedAt,&item.UpdatedAt,&item.Room.ID,&item.Room.RoomName)
		if err != nil{
			return reservations, err
		}
		reservations = append(reservations, item)
	}
	return reservations, nil
}


// AllNewReservations returns a slice of all NEW(processed=0) reservations
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservations, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservations

	query := ` select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
				r.end_date, r.room_id, r.created_at, r.updated_at,r.processed, rm.id, rm.room_name
				from reservations r
				left join rooms rm on (r.room_id = rm.id)
				where processed = 0
				order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()
	
	for rows.Next(){
		var item models.Reservations
		err := rows.Scan(&item.ID,&item.FirstName,&item.LastName,&item.Email,&item.Phone,&item.StartDate,
			&item.EndDate,&item.RoomID,&item.CreatedAt,&item.UpdatedAt,&item.Processed,&item.Room.ID,&item.Room.RoomName)
		if err != nil{
			return reservations, err
		}
		reservations = append(reservations, item)
	}
	return reservations, nil
}


// GetReservationByID return one reservation detail by ID
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservations, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservations

	query := `
		select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
		r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
		from reservations r
		left join rooms rm on (r.room_id = rm.id )
		where r.id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id )
	err := row.Scan(
		&res.ID,&res.FirstName,&res.LastName,&res.Email,&res.Phone,&res.StartDate,
		&res.EndDate,&res.RoomID,&res.CreatedAt,&res.UpdatedAt,&res.Processed,
		&res.Room.ID,&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}

	return res, nil

}

// UpdateReservation updates a reservation
func (m *postgresDBRepo) UpdateReservation(u models.Reservations, id int) (error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set first_name=$1, last_name=$2, email=$3, phone=$4, updated_at=$5
				where id = $6`
	
	_, err := m.DB.ExecContext(ctx, query,u.FirstName,u.LastName,u.Email,u.Phone,time.Now(),id)
	if err != nil{
		return err
	}

	return nil

}

// DeleteReservation delete a rerservation
func (m *postgresDBRepo) DeleteReservation(id int) (error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from reservations where id=$1`
	
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil{
		return err
	}

	return nil

}

// UpdateProcessedForReservation updates proceesed
func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) (error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update reservations set processed = $1 where id = $2`
	
	_, err := m.DB.ExecContext(ctx, query, processed, id)
	if err != nil{
		return err
	}

	return nil

}


// AllRooms gets all rooms
func (m *postgresDBRepo) AllRooms() ([]models.Room, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `select id, room_name, created_at, updated_at from rooms
				order by room_name`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil{
		return rooms, err
	}
	defer rows.Close()

	for rows.Next(){
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil{
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRestrictionsForRoomByDate get restrictions for rooms
func (m *postgresDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestrictions, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestrictions

	//coalesce(reservation_id, 0) : if reservation_id is NULL which means it's a block, replace it with 0.
	query := `select id, coalesce(reservation_id, 0), restriction_id, room_id, start_date, end_date
			from room_restrictions where $1<end_date and $2 >= start_date and room_id = $3`

	rows, err := m.DB.QueryContext(ctx, query,start, end, roomID)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		var r models.RoomRestrictions
		err := rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.RoomID,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil{
			return nil, err
		}
		restrictions = append(restrictions, r)
	}
	if err = rows.Err(); err!=nil{
		return nil, err
	}
	return restrictions, nil
}

// InsertBlockForRoom inserts blocks for a room
func (m *postgresDBRepo) InsertBlockForRoom(id int, startDate time.Time) error{
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into room_restrictions (start_date, end_date, room_id, restriction_id
		, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := m.DB.ExecContext(ctx, query, startDate, startDate.AddDate(0,0,1), id, 2, time.Now(), time.Now())
	if err!= nil {
		log.Println(err)
		return err
	}

	return nil
}

// DeleteBlockByID handles deleting blocks on calendar
func (m *postgresDBRepo) DeleteBlockByID(id int) error{
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from room_restrictions where id = $1`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err!= nil {
		log.Println(err)
		return err
	}

	return nil
}
