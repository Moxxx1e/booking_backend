package repository

import (
	"database/sql"
	"fmt"
	"github.com/booking_backend/internal/models"
	fixtureModels "github.com/booking_backend/internal/room/fixtures"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	sortPackage "sort"
	"testing"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Loader
)

func GetTestDBConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "booking_test")
}

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("postgres", GetTestDBConnString())
	if err != nil {
		log.Fatal(err)
	}

	fixtures, err = testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("../fixtures"),
	)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func prepareTestDatabase() {
	if err := fixtures.Load(); err != nil {
		log.Fatal(err)
	}
}

func TestRoomRepository_Insert_OK(t *testing.T) {
	prepareTestDatabase()

	roomRep := NewRoomRepository(db)
	roomModel := fixtureModels.NewDataBuilder().CreateNewRoomModel()

	err := roomRep.Insert(roomModel)

	// fixture id logic
	assert.NoError(t, err)
	assert.Equal(t, uint64(10001), roomModel.ID)
}

func TestRoomRepository_SelectByID(t *testing.T) {
	prepareTestDatabase()
	roomRep := NewRoomRepository(db)
	existedRoom := fixtureModels.NewDataBuilder().CreateFirstRoom()

	actualRoom, err := roomRep.SelectByID(existedRoom.ID)

	assert.NoError(t, err)
	assert.Equal(t, existedRoom, actualRoom)
}

func TestRoomRepository_SelectByID_NoRows(t *testing.T) {
	prepareTestDatabase()
	roomRep := NewRoomRepository(db)

	roomModel := fixtureModels.NewDataBuilder().CreateNewRoomModel()
	_, err := roomRep.SelectByID(roomModel.ID)

	assert.Error(t, err)
}

func TestRoomRepository_SelectRooms_Created_ASC(t *testing.T) {
	prepareTestDatabase()
	roomRep := NewRoomRepository(db)

	sort := &models.Sort{
		OrderBy: "created",
		Desc:    false,
	}

	existedRooms := fixtureModels.NewDataBuilder().CreateAllExistedRooms()
	sortPackage.Slice(existedRooms, func(i int, j int) bool {
		return existedRooms[i].Created.Before(existedRooms[j].Created)
	})

	actualRooms, err := roomRep.SelectRooms(sort)

	assert.NoError(t, err)
	assert.Equal(t, existedRooms, actualRooms)
}

func TestRoomRepository_SelectRooms_Price_ASC(t *testing.T) {
	prepareTestDatabase()
	roomRep := NewRoomRepository(db)

	sort := &models.Sort{
		OrderBy: "price",
		Desc:    false,
	}

	existedRooms := fixtureModels.NewDataBuilder().CreateAllExistedRooms()
	sortPackage.Slice(existedRooms, func(i int, j int) bool {
		return existedRooms[i].Price < existedRooms[j].Price
	})

	actualRooms, err := roomRep.SelectRooms(sort)

	assert.NoError(t, err)
	assert.Equal(t, existedRooms, actualRooms)
}

func TestRoomRepository_SelectRooms_Price_DESC(t *testing.T) {
	prepareTestDatabase()
	roomRep := NewRoomRepository(db)

	sort := &models.Sort{
		OrderBy: "price",
		Desc:    true,
	}

	existedRooms := fixtureModels.NewDataBuilder().CreateAllExistedRooms()
	sortPackage.Slice(existedRooms, func(i int, j int) bool {
		return existedRooms[i].Price > existedRooms[j].Price
	})

	actualRooms, err := roomRep.SelectRooms(sort)

	assert.NoError(t, err)
	assert.Equal(t, existedRooms, actualRooms)
}

func TestRoomRepository_SelectRooms_Created_DESC(t *testing.T) {
	prepareTestDatabase()
	roomRep := NewRoomRepository(db)

	sort := &models.Sort{
		OrderBy: "created",
		Desc:    true,
	}

	existedRooms := fixtureModels.NewDataBuilder().CreateAllExistedRooms()
	sortPackage.Slice(existedRooms, func(i int, j int) bool {
		return existedRooms[i].Created.After(existedRooms[j].Created)
	})

	actualRooms, err := roomRep.SelectRooms(sort)

	assert.NoError(t, err)
	assert.Equal(t, existedRooms, actualRooms)
}

func TestRoomRepository_DeleteRoomAndBookings(t *testing.T) {
	prepareTestDatabase()
	roomRep := NewRoomRepository(db)
	existedRoom := fixtureModels.NewDataBuilder().CreateFirstRoom()

	err := roomRep.DeleteRoomAndBookings(existedRoom.ID)

	assert.NoError(t, err)

	_, err = roomRep.SelectByID(existedRoom.ID)
	assert.Error(t, err)
}
