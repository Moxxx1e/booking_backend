package room

import "github.com/booking_backend/internal/models"

type RoomRepository interface {
	Insert(room *models.Room) error
	DeleteRoomAndBookings(id uint64) error
	SelectByID(id uint64) (*models.Room, error)
	SelectRooms(sort *models.Sort) ([]*models.Room, error)
}
