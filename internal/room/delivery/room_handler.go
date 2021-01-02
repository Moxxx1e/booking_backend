package delivery

import (
	. "github.com/booking_backend/internal/consts"
	"github.com/booking_backend/internal/helpers/errors"
	"github.com/booking_backend/internal/models"
	"github.com/booking_backend/internal/room"
	"github.com/booking_backend/tools/request_reader"
	"github.com/booking_backend/tools/response"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type RoomHandler struct {
	roomUseCase room.RoomUseCase
}

func NewRoomHandler(useCase room.RoomUseCase) *RoomHandler {
	return &RoomHandler{roomUseCase: useCase}
}

func (rh *RoomHandler) Configure(e *echo.Echo) {
	e.POST("rooms/create", rh.CreateRoom())
	e.GET("rooms/list", rh.GetRooms())
	e.DELETE("rooms/:id", rh.DeleteRoom())
}

func (rh *RoomHandler) CreateRoom() echo.HandlerFunc {
	type Request struct {
		// TODO: ограничить длину description
		// TODO: валидация цены
		Description string `form:"description" validate:"required"`
		Price       uint64 `form:"price" validate:"required"`
	}

	return func(context echo.Context) error {
		req := &Request{}
		if customErr := request_reader.NewRequestReader(context).Read(req); customErr != nil {
			logrus.Error(customErr)
			return context.JSON(customErr.HTTPCode, response.Response{Error: customErr})
		}

		room := &models.Room{
			Description: req.Description,
			Price:       req.Price,
			Created:     time.Now(),
		}

		if customErr := rh.roomUseCase.CreateRoom(room); customErr != nil {
			logrus.Error(customErr)
			return context.JSON(customErr.HTTPCode, response.Response{Error: customErr})
		}

		return context.JSON(http.StatusOK, response.Response{
			Body: &response.Body{
				"id": room.ID,
			},
		})
	}
}

func (rh *RoomHandler) GetRooms() echo.HandlerFunc {
	type Request struct {
		models.Sort
	}

	return func(context echo.Context) error {
		req := &Request{}
		if customErr := request_reader.NewRequestReader(context).Read(req); customErr != nil {
			logrus.Error(customErr)
			return context.JSON(customErr.HTTPCode, response.Response{Error: customErr})
		}

		if req.Sort.OrderBy == "" {
			req.Sort.OrderBy = "created"
		}

		rooms, customErr := rh.roomUseCase.GetRoomsList(&req.Sort)
		if customErr != nil {
			logrus.Error(customErr)
			return context.JSON(customErr.HTTPCode, response.Response{Error: customErr})
		}

		return context.JSON(http.StatusOK, response.Response{
			Body: &response.Body{
				"rooms": rooms,
			},
		})
	}
}

func (rh *RoomHandler) DeleteRoom() echo.HandlerFunc {
	return func(context echo.Context) error {
		roomID, parseErr := strconv.ParseUint(context.Param("id"), 10, 64)
		if parseErr != nil {
			customErr := errors.New(CodeInternalError, parseErr)
			logrus.Error(customErr)
			return context.JSON(customErr.HTTPCode, response.Response{Error: customErr})
		}

		customErr := rh.roomUseCase.DeleteRoomAndBookings(roomID)
		if customErr != nil {
			logrus.Error(customErr)
			return context.JSON(customErr.HTTPCode, response.Response{Error: customErr})
		}

		return context.JSON(http.StatusOK, response.Response{
			Message: "successfully deleted",
		})
	}
}
