package delivery

import (
	"github.com/booking_backend/internal/booking"
	. "github.com/booking_backend/internal/consts"
	"github.com/booking_backend/internal/helpers/errors"
	"github.com/booking_backend/internal/models"
	"github.com/booking_backend/tools/request_reader"
	"github.com/booking_backend/tools/response"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type BookingHandler struct {
	bookingUseCase booking.BookingUseCase
}

func NewBookingHandler(useCase booking.BookingUseCase) *BookingHandler {
	return &BookingHandler{bookingUseCase: useCase}
}

func (bh *BookingHandler) Configure(e *echo.Echo) {
	e.POST("bookings/create", bh.CreateBooking())
	e.GET("bookings/list", bh.GetRoomBookings())
	e.DELETE("bookings/:id", bh.DeleteBooking())
}

func (bh *BookingHandler) CreateBooking() echo.HandlerFunc {
	type Request struct {
		RoomID    uint64            `form:"room_id" validate:"required"`
		DateStart models.CustomDate `form:"date_start" validate:"required"`
		DateEnd   models.CustomDate `form:"date_end" validate:"required"`
	}

	return func(context echo.Context) error {
		req := &Request{}
		if customErr := request_reader.NewRequestReader(context).Read(req); customErr != nil {
			logrus.Info(customErr)
			return context.JSON(customErr.HTTPCode, response.Response{Error: customErr})
		}

		booking := &models.Booking{
			DateStart: req.DateStart.Date,
			DateEnd:   req.DateEnd.Date,
			Room:      req.RoomID,
		}

		if customErr := bh.bookingUseCase.CreateBooking(booking); customErr != nil {
			logrus.Info(customErr)
			return context.JSON(customErr.HTTPCode, response.Response{Error: customErr})
		}

		return context.JSON(http.StatusOK, response.Response{
			Body: &response.Body{
				"id": booking.ID,
			},
		})
	}
}

func (bh *BookingHandler) GetRoomBookings() echo.HandlerFunc {
	type Request struct {
		RoomID uint64 `query:"room_id"`
	}

	return func(context echo.Context) error {
		roomID, parseErr := strconv.ParseUint(context.QueryParam("room_id"), 10, 64)
		if parseErr != nil {
			customErr := errors.New(CodeInternalError, parseErr)
			logrus.Error(customErr)
			return context.JSON(customErr.HTTPCode, response.Response{Error: customErr})
		}

		bookings, customErr := bh.bookingUseCase.GetRoomBookings(roomID)
		if customErr != nil {
			logrus.Error(customErr)
			return context.JSON(customErr.HTTPCode, response.Response{Error: customErr})
		}

		return context.JSON(http.StatusOK, response.Response{
			Body: &response.Body{
				"bookings": bookings,
			},
		})
	}
}

func (bh *BookingHandler) DeleteBooking() echo.HandlerFunc {
	return func(context echo.Context) error {
		bookingID, parseErr := strconv.ParseUint(context.Param("id"), 10, 64)
		if parseErr != nil {
			customErr := errors.New(CodeInternalError, parseErr)
			logrus.Error(customErr)
			return context.JSON(customErr.HTTPCode, response.Response{Error: customErr})
		}

		customErr := bh.bookingUseCase.DeleteBooking(bookingID)
		if customErr != nil {
			logrus.Error(customErr)
			return context.JSON(customErr.HTTPCode, response.Response{Error: customErr})
		}

		return context.JSON(http.StatusOK, response.Response{
			Message: "successfully deleted",
		})
	}
}
