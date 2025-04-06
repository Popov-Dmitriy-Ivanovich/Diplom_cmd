package models

import "time"

type Event struct {
	Type      string
	TimeStamp time.Time
}
