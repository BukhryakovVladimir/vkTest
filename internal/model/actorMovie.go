package model

import "time"

type ActorMovie struct {
	MovieID   int       `json:"movieID"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Sex       string    `json:"sex"`
	BirthDate time.Time `json:"birthDate"`
}
