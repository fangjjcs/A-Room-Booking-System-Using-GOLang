package dbrepo

import (
	"context"
	"time"

	"github.com/fangjjcs/bookings-app/pkg/models"
)
func (m *postgresDBRepo) AllUsers() bool{
	return true
}

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

// Insert a room restriction into database
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


// Search Availability
func (m *postgresDBRepo) SearchAvailabilityByDate(start, end time.Time, roomID int) (bool,error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	stmt := `select count(id) from room_restrictions where
	        $1 < end_date and $2 > start_date and room_id = $3;`
	
	err := m.DB.QueryRowContext(ctx, stmt,start,end, roomID).Scan(&count)
	if err!=nil{
		return false, err
	}

	// assumes there exists only 1 room
	if count == 0 {
		return true, nil
	}
	return false, nil
}