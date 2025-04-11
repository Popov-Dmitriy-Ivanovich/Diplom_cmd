package models

type Event struct {
	ID        uint
	ActionID  uint
	Type      string
	TimeStamp TimeStamp
}
