package response

import "github.com/booking_backend/internal/helpers/errors"

type Body map[string]interface{}

type Response struct {
	Error   *errors.Error `json:"error,omitempty"`
	Message string        `json:"message,omitempty"`
	Body    *Body         `json:"body,omitempty"`
}
