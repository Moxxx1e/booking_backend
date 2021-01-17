package repository

import (
	"context"
	"database/sql"
	"github.com/booking_backend/internal/models"
	"github.com/booking_backend/internal/room"
	"github.com/sirupsen/logrus"
	"strings"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) room.RoomRepository {
	return &RoomRepository{db: db}
}

func (rep *RoomRepository) Insert(room *models.Room) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	err = tx.QueryRow(`
		INSERT INTO rooms(description, price, created) 
		VALUES ($1, $2, $3) RETURNING id`,
		room.Description, room.Price, room.Created).
		Scan(&room.ID)
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

func (rep *RoomRepository) SelectByID(id uint64) (*models.Room, error) {
	room := &models.Room{}
	err := rep.db.QueryRow(`
		SELECT id, description, price, created
		FROM rooms
		WHERE id=$1`, id).
		Scan(&room.ID, &room.Description, &room.Price, &room.Created)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (rep *RoomRepository) DeleteRoomAndBookings(id uint64) error {
	tx, err := rep.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	// Bookings will be deleted cascade
	_, err = tx.Exec(`
		DELETE
		FROM rooms
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

func createSelectQuery(sort *models.Sort) string {
	query := "SELECT id, description, price, created FROM rooms"
	switch sort.OrderBy {
	case "price":
		query = strings.Join([]string{query, "ORDER BY price"}, " ")
	case "created":
		query = strings.Join([]string{query, "ORDER BY created"}, " ")
	}
	if sort.Desc {
		query = strings.Join([]string{query, "DESC"}, " ")
	}
	return query
}

func (rep *RoomRepository) SelectRooms(sort *models.Sort) ([]*models.Room, error) {
	query := createSelectQuery(sort)

	rows, err := rep.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		room := &models.Room{}
		if err := rows.Scan(&room.ID, &room.Description,
			&room.Price, &room.Created); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}
