package models

import "time"

type CustomDate struct {
	Date string
}

func (ct *CustomDate) UnmarshalParam(param string) error {
	t, err := time.Parse(`2006-01-02`, param)
	if err != nil {
		return err
	}
	ct.Date = t.Format(`2006-01-02`)
	return nil
}

type Booking struct {
	ID        uint64 `json:"booking_id"`
	DateStart string `json:"date_start"`
	DateEnd   string `json:"date_end"`
	Room      uint64 `json:"room"`
}
