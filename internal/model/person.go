package model

import "time"

type Person struct {
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Sex       string    `json:"sex"`
	BirthDate time.Time `json:"birthDate"`
}
