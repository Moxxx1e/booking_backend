package main

import (
	"database/sql"
	"fmt"
	bookingDelivery "github.com/booking_backend/internal/booking/delivery"
	bookingRepository "github.com/booking_backend/internal/booking/repository"
	bookingUseCase "github.com/booking_backend/internal/booking/usecases"
	roomDelivery "github.com/booking_backend/internal/room/delivery"
	roomRepository "github.com/booking_backend/internal/room/repository"
	roomUseCase "github.com/booking_backend/internal/room/usecases"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
)

const serverAddr = "localhost:9000"

func GetDbConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "booking")
}

func main() {
	// Database
	dbConnection, err := sql.Open("postgres", GetDbConnString())
	if err != nil {
		log.Fatal(err)
	}
	defer dbConnection.Close()

	if err := dbConnection.Ping(); err != nil {
		log.Fatal(err)
	}

	roomRepo := roomRepository.NewRoomRepository(dbConnection)
	roomUseCase := roomUseCase.NewRoomUseCase(roomRepo)
	roomHandler := roomDelivery.NewRoomHandler(roomUseCase)

	bookingRepo := bookingRepository.NewBookingRepository(dbConnection)
	bookingUseCase := bookingUseCase.NewBookingUseCase(bookingRepo, roomRepo)
	bookingHandler := bookingDelivery.NewBookingHandler(bookingUseCase)

	e := echo.New()

	roomHandler.Configure(e)
	bookingHandler.Configure(e)

	log.Fatal(e.Start(serverAddr))
}
