package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/booking_backend/internal/booking/mocks"
	"github.com/booking_backend/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var bookingModel = &models.Booking{
	ID:        1,
	DateStart: "2020-12-10",
	DateEnd:   "2021-12-10",
	Room:      4,
}

var firstRoom = &models.Room{
	ID: 1,
}

var bookingsOfFirstRoom = []*models.Booking{
	&models.Booking{
		ID:        1,
		DateStart: "2020-12-10",
		DateEnd:   "2021-12-10",
		Room:      1,
	},
	&models.Booking{
		ID:        5,
		DateStart: "2020-12-10",
		DateEnd:   "2021-12-10",
		Room:      1,
	},
	&models.Booking{
		ID:        7,
		DateStart: "2020-12-10",
		DateEnd:   "2021-12-10",
		Room:      1,
	}, &models.Booking{
		ID:        97,
		DateStart: "2020-12-10",
		DateEnd:   "2021-12-10",
		Room:      1,
	},
}

func TestBookingRepository_Insert(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	bookingPgRep := NewBookingRepository(db)

	mocks.MockInsertSuccess(mock, bookingModel)
	err = bookingPgRep.Insert(bookingModel)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBookingRepository_SelectByID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	bookingPgRep := NewBookingRepository(db)

	mocks.MockSelectReturnRows(mock, bookingModel)
	resultBooking, err := bookingPgRep.SelectByID(bookingModel.ID)

	assert.NoError(t, err)
	assert.Equal(t, bookingModel, resultBooking)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBookingRepository_SelectByID_ErrNoRows(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	bookingPgRep := NewBookingRepository(db)

	mocks.MockSelectBookingByIDReturnErrNoRows(mock, bookingModel.ID)
	resultBooking, err := bookingPgRep.SelectByID(bookingModel.ID)

	assert.Error(t, sql.ErrNoRows)
	assert.Nil(t, resultBooking)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBookingRepository_SelectRoomBookings(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	bookingPgRep := NewBookingRepository(db)

	mocks.MockSelectBookingList(mock, firstRoom.ID, bookingsOfFirstRoom)
	resultBooking, err := bookingPgRep.SelectRoomBookings(firstRoom.ID)

	assert.NoError(t, err)
	assert.Equal(t, bookingsOfFirstRoom, resultBooking)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBookingRepository_SelectRoomBookings_Nil(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	bookingPgRep := NewBookingRepository(db)

	mocks.MockSelectBookingList(mock, firstRoom.ID, nil)
	resultBooking, err := bookingPgRep.SelectRoomBookings(firstRoom.ID)

	assert.NoError(t, err)
	assert.Nil(t, resultBooking)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBookingRepository_DeleteByID(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	bookingPgRep := NewBookingRepository(db)

	mocks.MockDeleteSuccess(mock, bookingModel.ID)
	err = bookingPgRep.DeleteByID(bookingModel.ID)

	assert.NoError(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
