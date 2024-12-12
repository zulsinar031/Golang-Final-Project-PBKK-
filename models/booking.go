package models

import "time"

type Booking struct {
	ID            int
	Arrivaldate   time.Time
	Departuredate time.Time
	Hotelname     string
	Username      string
	Comment       string
}