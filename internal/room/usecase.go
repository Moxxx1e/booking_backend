package room

import (
	"github.com/booking_backend/internal/helpers/errors"
	"github.com/booking_backend/internal/models"
)

type RoomUseCase interface {
	CreateRoom(room *models.Room) *errors.Error
	DeleteRoomAndBookings(id uint64) *errors.Error
	GetRoomsList(sort *models.Sort) ([]*models.Room, *errors.Error)
}
