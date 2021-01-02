package booking

import "github.com/booking_backend/internal/models"

type BookingRepository interface {
	Insert(booking *models.Booking) error
	SelectByID(id uint64) (*models.Booking, error)
	DeleteByID(id uint64) error
	SelectRoomBookings(roomID uint64) ([]*models.Booking, error)
}
