package fixtures

import (
	"github.com/booking_backend/internal/models"
	"log"
	sortPackage "sort"
	"time"
)

type DataBuilder struct{}

func NewDataBuilder() *DataBuilder{
	return &DataBuilder{}
}

func (db *DataBuilder) CreateNewRoomModel() *models.Room {
	return &models.Room{
		ID:          0,
		Description: "Just a new room",
		Price:       100,
		Created:     time.Now(),
	}
}

func(db *DataBuilder)  CreateFirstRoom() *models.Room {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}
	var existedRoom = &models.Room{
		ID:          1,
		Description: "room at the Hotel California",
		Price:       500,
		Created:     time.Date(2021, 1, 8, 19, 37, 51, 0, loc),
	}
	return existedRoom
}

func (db *DataBuilder) CreateRoomsWithoutForthOrderByCreate() []*models.Room {
	rooms := db.CreateAllExistedRooms()
	rooms = append(rooms[:3], rooms[4:]...)
	sortPackage.Slice(rooms, func(i int, j int) bool {
		return rooms[i].Created.Before(rooms[j].Created)
	})
	return rooms
}

func (db *DataBuilder) CreateAllExistedRooms() []*models.Room {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal(err)
	}
	existedRooms := []*models.Room{
		&models.Room{
			ID:          1,
			Description: "room at the Hotel California",
			Price:       500,
			Created:     time.Date(2021, 1, 8, 19, 37, 51, 0, loc),
		},
		&models.Room{
			ID:          2,
			Description: "room at the Grand Budapest Hotel",
			Price:       11500,
			Created:     time.Date(2021, 1, 9, 19, 37, 51, 0, loc),
		}, &models.Room{
			ID:          3,
			Description: "room at the Hostel Teriba",
			Price:       750,
			Created:     time.Date(2021, 1, 7, 19, 37, 51, 0, loc),
		}, &models.Room{
			ID:          4,
			Description: "room at the Hostel Friends",
			Price:       300,
			Created:     time.Date(2021, 1, 6, 19, 37, 51, 0, loc),
		},
	}
	return existedRooms
}
