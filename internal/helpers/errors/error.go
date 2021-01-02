package errors

import (
	. "github.com/booking_backend/internal/consts"
	"net/http"
)

type Error struct {
	Code        uint64 `json:"code"`
	HTTPCode    int    `json:"-"`
	Message     string `json:"message"`
	UserMessage string `json:"user_message"`
}

var WrongErrorCode = &Error{
	HTTPCode:    http.StatusTeapot,
	Message:     "wrong error code",
	UserMessage: "Что-то пошло не так",
}

func New(code uint64, err error) *Error {
	customErr, has := Errors[code]
	if !has {
		return WrongErrorCode
	}
	customErr.Message = err.Error()
	return customErr
}

func Get(code uint64) *Error {
	err, has := Errors[code]
	if !has {
		return WrongErrorCode
	}
	return err
}

var Errors = map[uint64]*Error{
	CodeInternalError: {
		Code:        CodeInternalError,
		HTTPCode:    http.StatusInternalServerError,
		Message:     "something went wrong",
		UserMessage: "Что-то пошло не так",
	},
	CodeBadRequest: {
		Code:        CodeBadRequest,
		HTTPCode:    http.StatusBadRequest,
		Message:     "wrong request data",
		UserMessage: "Неверный формат запроса",
	},
	CodeRoomDoesNotExist: {
		Code:        CodeRoomDoesNotExist,
		HTTPCode:    http.StatusNotFound,
		Message:     "room with this id doesn't exist",
		UserMessage: "Комнаты с таким ID не существует",
	},
}
