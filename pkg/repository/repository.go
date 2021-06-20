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
	SearchAvailabilityByDate(start, end time.Time, roomID int) (bool,error)
}