package usecases

import (
	"database/sql"
	"github.com/booking_backend/internal/booking/mocks"
	"github.com/booking_backend/internal/consts"
	"github.com/booking_backend/internal/helpers/errors"
	"github.com/booking_backend/internal/models"
	mockRoom "github.com/booking_backend/internal/room/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var bookingModel = &models.Booking{
	ID:        3,
	DateStart: "2022-01-02",
	DateEnd:   "2022-01-02",
	Room:      1,
}

var firstRoom = &models.Room{
	ID:          1,
	Description: "some description",
	Price:       500,
	Created:     time.Time{},
}

var bookings = []*models.Booking{
	&models.Booking{
		ID:        1,
		DateStart: "2022-01-02",
		DateEnd:   "2023-01-02",
		Room:      1,
	},
	&models.Booking{
		ID:        2,
		DateStart: "2022-01-02",
		DateEnd:   "2023-01-02",
		Room:      1,
	},
	&models.Booking{
		ID:        4,
		DateStart: "2022-01-02",
		DateEnd:   "2023-01-02",
		Room:      1,
	},
}

func TestBookingUseCase_CreateBooking_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookingRep := mocks.NewMockBookingRepository(ctrl)
	roomRep := mockRoom.NewMockRoomRepository(ctrl)
	bookingUseCase := NewBookingUseCase(bookingRep, roomRep)

	roomRep.
		EXPECT().
		SelectByID(bookingModel.Room).
		Return(firstRoom, nil)

	bookingRep.
		EXPECT().
		Insert(bookingModel).
		Return(nil)

	err := bookingUseCase.CreateBooking(bookingModel)
	assert.Equal(t, (*errors.Error)(nil), err)
}

func TestBookingUseCase_CreateBooking_RoomDoesNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookingRep := mocks.NewMockBookingRepository(ctrl)
	roomRep := mockRoom.NewMockRoomRepository(ctrl)
	bookingUseCase := NewBookingUseCase(bookingRep, roomRep)

	roomRep.
		EXPECT().
		SelectByID(bookingModel.Room).
		Return(nil, sql.ErrNoRows)

	err := bookingUseCase.CreateBooking(bookingModel)
	assert.Equal(t, errors.Get(consts.CodeRoomDoesNotExist), err)
}

func TestBookingUseCase_GetRoomBookings_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookingRep := mocks.NewMockBookingRepository(ctrl)
	roomRep := mockRoom.NewMockRoomRepository(ctrl)
	bookingUseCase := NewBookingUseCase(bookingRep, roomRep)

	roomRep.
		EXPECT().
		SelectByID(bookingModel.Room).
		Return(firstRoom, nil)

	bookingRep.
		EXPECT().
		SelectRoomBookings(bookingModel.Room).
		Return(bookings, nil)

	bookingsResult, err := bookingUseCase.GetRoomBookings(bookingModel.Room)
	assert.Equal(t, (*errors.Error)(nil), err)
	assert.Equal(t, bookings, bookingsResult)
}

func TestBookingUseCase_GetRoomBookings_NoBookings(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookingRep := mocks.NewMockBookingRepository(ctrl)
	roomRep := mockRoom.NewMockRoomRepository(ctrl)
	bookingUseCase := NewBookingUseCase(bookingRep, roomRep)

	roomRep.
		EXPECT().
		SelectByID(bookingModel.Room).
		Return(firstRoom, nil)

	bookingRep.
		EXPECT().
		SelectRoomBookings(bookingModel.Room).
		Return(nil, nil)

	bookings, err := bookingUseCase.GetRoomBookings(bookingModel.Room)
	assert.Equal(t, (*errors.Error)(nil), err)
	assert.Equal(t, []*models.Booking{}, bookings)
}

func TestBookingUseCase_GetRoomBookings_RoomDoesNotExist(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookingRep := mocks.NewMockBookingRepository(ctrl)
	roomRep := mockRoom.NewMockRoomRepository(ctrl)
	bookingUseCase := NewBookingUseCase(bookingRep, roomRep)

	roomRep.
		EXPECT().
		SelectByID(bookingModel.Room).
		Return(nil, sql.ErrNoRows)

	bookings, err := bookingUseCase.GetRoomBookings(bookingModel.Room)
	assert.Equal(t, errors.Get(consts.CodeRoomDoesNotExist), err)
	assert.Nil(t, bookings)
}

func TestBookingUseCase_DeleteBooking_NoBooking(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookingRep := mocks.NewMockBookingRepository(ctrl)
	roomRep := mockRoom.NewMockRoomRepository(ctrl)
	bookingUseCase := NewBookingUseCase(bookingRep, roomRep)

	bookingRep.
		EXPECT().
		SelectByID(bookingModel.ID).
		Return(nil, sql.ErrNoRows)

	err := bookingUseCase.DeleteBooking(bookingModel.ID)
	assert.Equal(t, errors.Get(consts.CodeBookingDoesNotExist), err)
}

func TestBookingUseCase_DeleteBooking_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bookingRep := mocks.NewMockBookingRepository(ctrl)
	roomRep := mockRoom.NewMockRoomRepository(ctrl)
	bookingUseCase := NewBookingUseCase(bookingRep, roomRep)

	bookingRep.
		EXPECT().
		SelectByID(bookingModel.ID).
		Return(bookingModel, nil)
	bookingRep.
		EXPECT().
		DeleteByID(bookingModel.ID).
		Return(nil)

	err := bookingUseCase.DeleteBooking(bookingModel.ID)
	assert.Equal(t, (*errors.Error)(nil), err)
}
