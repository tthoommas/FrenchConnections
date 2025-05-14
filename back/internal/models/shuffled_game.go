package models

import "time"

type ShuffledGame struct {
	CreatedBy     string    `json:"createdBy"`
	CreationDate  time.Time `json:"creationDate"`
	ShuffledWords []string  `json:"game"`
}
