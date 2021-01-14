package mocks

import (
	"database/sql"
	"github.com/booking_backend/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func MockInsertSuccess(mock sqlmock.Sqlmock, booking *models.Booking) {
	mock.ExpectBegin()
	rows := sqlmock.NewRows([]string{"id"}).AddRow(booking.ID)
	mock.ExpectQuery(`INSERT INTO bookings`).
		WithArgs(booking.DateStart, booking.DateEnd, booking.Room).
		WillReturnRows(rows)
	mock.ExpectCommit()
}

func MockDeleteSuccess(mock sqlmock.Sqlmock, id uint64) {
	mock.ExpectBegin()
	res := sqlmock.NewResult(0, 1)
	mock.ExpectExec(`DELETE FROM bookings`).
		WithArgs(id).
		WillReturnResult(res)
	mock.ExpectCommit()
}

func MockSelectBookingByIDReturnErrNoRows(mock sqlmock.Sqlmock, id uint64) {
	mock.ExpectQuery(`SELECT`).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)
}

func MockSelectReturnRows(mock sqlmock.Sqlmock, booking *models.Booking) {
	rows := sqlmock.NewRows([]string{"id", "date_start", "date_end", "room"})
	rows.AddRow(booking.ID, booking.DateStart, booking.DateEnd, booking.Room)
	mock.ExpectQuery(`SELECT`).
		WithArgs(booking.ID).
		WillReturnRows(rows)
}

func MockSelectBookingList(mock sqlmock.Sqlmock,
	roomID uint64, resultBookings []*models.Booking) {
	rows := sqlmock.NewRows([]string{"id", "date_start", "date_end", "room"})
	for _, booking := range resultBookings {
		rows.AddRow(booking.ID, booking.DateStart,
			booking.DateEnd, booking.Room)
	}

	mock.ExpectQuery(`SELECT`).
		WithArgs(roomID).
		WillReturnRows(rows)
}
