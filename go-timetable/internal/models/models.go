package models

import (
	"time"
	"github.com/gofrs/uuid"
)


type User struct {
	Id   *uuid.UUID `json:"id"`
	Name string     `json:"name"`
}

type Event struct {
	Id               string    `json:"id"`
	Summary          string    `json:"summary"`
	Location         *string   `json:"location,omitempty"`
	Description      *string   `json:"description,omitempty"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	Created          time.Time `json:"created"`
	LastModifiedTime time.Time `json:"lastModifiedTime"`
	Sequence         int       `json:"sequence"`
}
