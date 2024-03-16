package model

import "time"

type SearchMovie struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Date           time.Time `json:"date"`
	Rating         int16     `json:"rating"`
	ActorFirstName string    `json:"actorFirstName"`
	ActorLastName  string    `json:"actorLastName"`
}
