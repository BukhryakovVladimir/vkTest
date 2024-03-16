package model

import "time"

type ActorAndMovies struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Sex       string    `json:"sex"`
	BirthDate time.Time `json:"birthDate"`
	Movies    []Movie   `json:"movies"`
}
