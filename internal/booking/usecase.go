package booking

import (
	"github.com/booking_backend/internal/helpers/errors"
	"github.com/booking_backend/internal/models"
)

type BookingUseCase interface {
	CreateBooking(booking *models.Booking) *errors.Error
	DeleteBooking(id uint64) *errors.Error
	GetRoomBookings(roomID uint64) ([]*models.Booking, *errors.Error)
}
