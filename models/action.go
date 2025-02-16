package models

import "time"

type Action struct {
	ID uint `json:"-"`
	CreatedAt time.Time `json:"-"`

	Status ActionStatus `json:"-"`
	StatusID uint `json:"-"`
	
	LastLaunch *DateOnly

	Name string
	ShortDesc string
	Description string
	Cmd string

}