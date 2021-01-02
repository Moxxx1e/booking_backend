package usecases

import (
	"database/sql"
	"github.com/booking_backend/internal/booking"
	"github.com/booking_backend/internal/consts"
	"github.com/booking_backend/internal/helpers/errors"
	"github.com/booking_backend/internal/models"
	"github.com/booking_backend/internal/room"
	"time"
)

func NewBookingUseCase(bookingRepository booking.BookingRepository,
	roomRepository room.RoomRepository) booking.BookingUseCase {
	return &BookingUseCase{bookingRepo: bookingRepository,
		roomRepo: roomRepository}
}

type BookingUseCase struct {
	bookingRepo booking.BookingRepository
	roomRepo    room.RoomRepository
}

func checkDates(booking *models.Booking) *errors.Error {
	dateStart, err := time.Parse(`2006-01-02`, booking.DateStart)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	dateEnd, err := time.Parse(`2006-01-02`, booking.DateEnd)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	if dateEnd.Before(dateStart) {
		return errors.Get(consts.CodeIncorrectDates)
	}
	return nil
}

func (uc *BookingUseCase) CreateBooking(booking *models.Booking) *errors.Error {
	if err := checkDates(booking); err != nil {
		return err
	}

	_, err := uc.roomRepo.SelectByID(booking.Room)
	if err == sql.ErrNoRows {
		return errors.Get(consts.CodeRoomDoesNotExist)
	} else if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}

	err = uc.bookingRepo.Insert(booking)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *BookingUseCase) DeleteBooking(id uint64) *errors.Error {
	_, err := uc.bookingRepo.SelectByID(id)
	if err == sql.ErrNoRows {
		return errors.Get(consts.CodeBookingDoesNotExist)
	} else if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}

	err = uc.bookingRepo.DeleteByID(id)
	if err != nil {
		return errors.New(consts.CodeInternalError, err)
	}
	return nil
}

func (uc *BookingUseCase) GetRoomBookings(roomID uint64) ([]*models.Booking, *errors.Error) {
	_, err := uc.roomRepo.SelectByID(roomID)
	if err == sql.ErrNoRows {
		return nil, errors.Get(consts.CodeRoomDoesNotExist)
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}

	bookings, err := uc.bookingRepo.SelectRoomBookings(roomID)
	if bookings == nil && err == nil {
		return []*models.Booking{}, nil
	} else if err != nil {
		return nil, errors.New(consts.CodeInternalError, err)
	}
	return bookings, nil
}
