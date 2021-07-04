package repository

import (
	"time"

	"github.com/fangjjcs/bookings-app/pkg/models"
)

// Interface for different demand of database type
type DatabaseRepo interface{
	AllUsers() bool
	InsertReservations(res models.Reservations) (int,error)
	InsertRoomRestriction(r models.RoomRestrictions) error
	SearchAvailabilityByDatesAndRoomID(start, end time.Time, roomID int) (bool,error)
	SearchAvailibilityForAllRooms(start, end time.Time) ([]models.Room, error)

	GetUserByID(id int) (models.User, error)
	Authenticate(email, testPassword string) (int, string, error)
}