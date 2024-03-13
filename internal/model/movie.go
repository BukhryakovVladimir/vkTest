package model

import "time"

type Movie struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Rating      int16     `json:"rating"`
}
