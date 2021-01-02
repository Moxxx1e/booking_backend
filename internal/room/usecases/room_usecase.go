package usecases

import (
	"database/sql"
	"github.com/booking_backend/internal/consts"
	"github.com/booking_backend/internal/helpers/errors"
	"github.com/booking_backend/internal/models"
	"github.com/booking_backend/internal/room"
)

type RoomUseCase struct {
	roomsRep room.RoomRepository
}

func NewRoomUseCase(rep room.RoomRepository) room.RoomUseCase {
	return &RoomUseCase{roomsRep: rep}
}

func (uc *RoomUseCase) CreateRoom(room *models.Room) *errors.Error {
	err := uc.roomsRep.Insert(room)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *RoomUseCase) DeleteRoomAndBookings(id uint64) *errors.Error {
	_, err := uc.roomsRep.SelectByID(id)
	if err == sql.ErrNoRows {
		return errors.Get(consts.CodeRoomDoesNotExist)
	} else if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}

	err = uc.roomsRep.DeleteRoomAndBookings(id)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *RoomUseCase) GetRoomsList(sort *models.Sort) ([]*models.Room, *errors.Error) {
	rooms, err := uc.roomsRep.SelectRooms(sort)
	if err == nil && rooms == nil {
		return []*models.Room{}, nil
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	return rooms, nil
}
