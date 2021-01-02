package repository

import (
	"context"
	"database/sql"
	"github.com/booking_backend/internal/booking"
	"github.com/booking_backend/internal/models"
	"github.com/sirupsen/logrus"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) booking.BookingRepository {
	return &BookingRepository{db: db}
}

func (rep *BookingRepository) Insert(booking *models.Booking) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = tx.QueryRow(`
		INSERT INTO bookings(date_start, date_end, room) 
		VALUES ($1, $2, $3) RETURNING id`,
		booking.DateStart, booking.DateEnd, booking.Room).
		Scan(&booking.ID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Info(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (rep *BookingRepository) SelectByID(id uint64) (*models.Booking, error) {
	booking := &models.Booking{}
	err := rep.db.QueryRow(`
		SELECT id, date_start, date_end, room
		FROM bookings
		WHERE id=$1`, id).
		Scan(&booking.ID, &booking.DateStart, &booking.DateEnd, &booking.Room)
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (rep *BookingRepository) DeleteByID(id uint64) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE
		FROM bookings
		WHERE id=$1`, id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			logrus.Info(rollbackErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (rep *BookingRepository) SelectRoomBookings(roomID uint64) ([]*models.Booking, error) {
	rows, err := rep.db.Query(`
		SELECT id, date_start, date_end, room
		FROM bookings
		WHERE room=$1
		ORDER BY date_start`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		booking := &models.Booking{}
		if err := rows.Scan(&booking.ID, &booking.DateStart,
			&booking.DateEnd, &booking.Room); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return bookings, nil
}
