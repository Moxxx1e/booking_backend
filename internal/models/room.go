package models

import "time"

type Room struct {
	ID          uint64    `json:"room_id"`
	Description string    `json:"description"`
	Price       uint64    `json:"price"`
	Created     time.Time `json:"created"`
}
