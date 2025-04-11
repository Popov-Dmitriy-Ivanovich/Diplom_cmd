package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// DateOnly
// Структура, сохраняющаяся в БД и представляющая дату без времени
type TimeStamp struct {
	time.Time
}

func (do TimeStamp) Value() (driver.Value, error) {
	return do.Add(0), nil
}

func (do *TimeStamp) Scan(src any) error {
	time, ok := src.(time.Time)
	if !ok {
		return errors.New("not time provided to DateOnly")
	}
	*do = TimeStamp{time}
	return nil
}

func (do TimeStamp) MarshalJSON() ([]byte, error) {
	timeStr := do.Format(time.RFC3339)
	return json.Marshal(timeStr)
}

func (do *TimeStamp) UnmarshalJSON(data []byte) (err error) {
	do.Time, err = time.Parse(time.DateOnly, string(data[1:len(data)-1]))
	if err != nil {
		do.Time, err = time.Parse(time.RFC3339, string(data[1:len(data)-1]))
	}
	return
}

func (do *TimeStamp) ToTime() time.Time {
	return do.Add(0)
}
