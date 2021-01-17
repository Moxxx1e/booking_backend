package usecases

import (
	"database/sql"
	"fmt"
	bookingRepository "github.com/booking_backend/internal/booking/repository"
	bookingUseCase "github.com/booking_backend/internal/booking/usecases"
	"github.com/booking_backend/internal/consts"
	"github.com/booking_backend/internal/helpers/errors"
	"github.com/booking_backend/internal/models"
	fixtureModels "github.com/booking_backend/internal/room/fixtures"
	"github.com/booking_backend/internal/room/repository"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"sort"
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

func TestRoomUseCase_CreateRoom_OK(t *testing.T) {
	prepareTestDatabase()

	rep := repository.NewRoomRepository(db)
	roomUseCase := NewRoomUseCase(rep)
	roomModel := fixtureModels.NewDataBuilder().CreateNewRoomModel()

	err := roomUseCase.CreateRoom(roomModel)

	assert.Nil(t, err)
	assert.Equal(t, uint64(10001), roomModel.ID)
}

func TestRoomUseCase_DeleteRoomAndBookings(t *testing.T) {
	prepareTestDatabase()
	roomRepository := repository.NewRoomRepository(db)
	roomUseCase := NewRoomUseCase(roomRepository)
	bookingRep := bookingRepository.NewBookingRepository(db)
	bookingUseCase := bookingUseCase.NewBookingUseCase(bookingRep, roomRepository)

	customErr := roomUseCase.DeleteRoomAndBookings(4)
	assert.Nil(t, customErr)

	rooms, customErr := roomUseCase.GetRoomsList(&models.Sort{
		OrderBy: "created",
		Desc:    false,
	})
	assert.Nil(t, customErr)
	assert.Equal(t, fixtureModels.NewDataBuilder().CreateRoomsWithoutForthOrderByCreate(), rooms)

	_, customErr = bookingUseCase.GetRoomBookings(4)
	assert.Equal(t, errors.Get(consts.CodeRoomDoesNotExist), customErr)

	bookings, err := bookingRep.SelectRoomBookings(4)
	assert.NoError(t, err)
	assert.Nil(t, bookings)
}

func TestRoomUseCase_GetRoomsList(t *testing.T) {
	prepareTestDatabase()
	roomRepository := repository.NewRoomRepository(db)
	roomUseCase := NewRoomUseCase(roomRepository)

	rooms, customErr := roomUseCase.GetRoomsList(&models.Sort{
		OrderBy: "created",
		Desc:    false,
	})

	expectedRooms := fixtureModels.NewDataBuilder().CreateAllExistedRooms()
	sort.Slice(expectedRooms, func(i, j int) bool {
		return expectedRooms[i].Created.Before(expectedRooms[j].Created)
	})

	assert.Nil(t, customErr)
	assert.Equal(t, expectedRooms, rooms)
}

func TestRoomUseCase_GetRoomsList_Created_ASC(t *testing.T) {
	prepareTestDatabase()
	roomRepository := repository.NewRoomRepository(db)
	roomUseCase := NewRoomUseCase(roomRepository)

	rooms, customErr := roomUseCase.GetRoomsList(&models.Sort{
		OrderBy: "created",
		Desc:    true,
	})

	expectedRooms := fixtureModels.NewDataBuilder().CreateAllExistedRooms()
	sort.Slice(expectedRooms, func(i, j int) bool {
		return expectedRooms[i].Created.After(expectedRooms[j].Created)
	})

	assert.Nil(t, customErr)
	assert.Equal(t, expectedRooms, rooms)
}

func TestRoomUseCase_GetRoomsList_Price(t *testing.T) {
	prepareTestDatabase()
	roomRepository := repository.NewRoomRepository(db)
	roomUseCase := NewRoomUseCase(roomRepository)

	rooms, customErr := roomUseCase.GetRoomsList(&models.Sort{
		OrderBy: "price",
		Desc:    false,
	})

	expectedRooms := fixtureModels.NewDataBuilder().CreateAllExistedRooms()
	sort.Slice(expectedRooms, func(i, j int) bool {
		return expectedRooms[i].Price < expectedRooms[j].Price
	})

	assert.Nil(t, customErr)
	assert.Equal(t, expectedRooms, rooms)
}

func TestRoomUseCase_GetRoomsList_Price_DESC(t *testing.T) {
	prepareTestDatabase()
	roomRepository := repository.NewRoomRepository(db)
	roomUseCase := NewRoomUseCase(roomRepository)

	rooms, customErr := roomUseCase.GetRoomsList(&models.Sort{
		OrderBy: "price",
		Desc:    true,
	})

	expectedRooms := fixtureModels.NewDataBuilder().CreateAllExistedRooms()
	sort.Slice(expectedRooms, func(i, j int) bool {
		return expectedRooms[i].Price > expectedRooms[j].Price
	})

	assert.Nil(t, customErr)
	assert.Equal(t, expectedRooms, rooms)
}

func TestRoomUseCase_GetRoomsList_Empty(t *testing.T) {
	prepareTestDatabase()
	roomRepository := repository.NewRoomRepository(db)
	roomUseCase := NewRoomUseCase(roomRepository)

	for _, id := range []uint64{1, 2, 3, 4} {
		customErr := roomUseCase.DeleteRoomAndBookings(id)
		assert.Nil(t, customErr)
	}

	rooms, customErr := roomUseCase.GetRoomsList(&models.Sort{
		OrderBy: "price",
		Desc:    true,
	})

	assert.Nil(t, customErr)
	assert.Equal(t, []*models.Room{}, rooms)
}
