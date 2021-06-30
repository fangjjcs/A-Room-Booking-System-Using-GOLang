package dbrepo

import (
	"context"
	"log"
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
func (m *postgresDBRepo) SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool,error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println(start,end, roomID)
	var count int
	stmt := `select count(id) from room_restrictions where
	        $1 < end_date and $2 > start_date and room_id = $3;`
	
	err := m.DB.QueryRowContext(ctx, stmt , start, end, roomID).Scan(&count)
	if err!=nil{
		return false, err
	}

	log.Println(count)
	// assumes there exists only 1 room
	if count == 0 {
		return true, nil
	}
	return false, nil
}


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