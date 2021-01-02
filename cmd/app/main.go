package main

import (
	"database/sql"
	"fmt"
	"github.com/booking_backend/internal/room/delivery"
	"github.com/booking_backend/internal/room/repository"
	"github.com/booking_backend/internal/room/usecases"
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

	// Repository
	roomRepo := repository.NewRoomRepository(dbConnection)
	roomUseCase := usecases.NewRoomUseCase(roomRepo)
	roomHandler := delivery.NewRoomHandler(roomUseCase)

	e := echo.New()

	roomHandler.Configure(e)

	log.Fatal(e.Start(serverAddr))
}
