package model

import "time"

type Movie struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Rating      int16     `json:"rating"`
	Actors      []Actor   `json:"actors,omitempty"`
}
