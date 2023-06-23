package repository

import (
	"github.com/zongjie233/udemy_lesson/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(r models.RoomRestriction) error

	SearchAvailabilityByDates(start, end time.Time, roomID int) (bool, error)
}
