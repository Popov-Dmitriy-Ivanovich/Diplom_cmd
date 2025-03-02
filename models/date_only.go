package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// DateOnly
// Структура, сохраняющаяся в БД и представляющая дату без времени
type DateOnly struct {
	time.Time
}

func (do DateOnly) Value() (driver.Value, error) {
	return do.Add(0), nil
}

func (do *DateOnly) Scan(src any) error {
	time, ok := src.(time.Time)
	if !ok {
		return errors.New("not time provided to DateOnly")
	}
	*do = DateOnly{time}
	return nil
}

func (do DateOnly) MarshalJSON() ([]byte, error) {
	timeStr := do.Format("02.01.06")
	return json.Marshal(timeStr)
}

func (do *DateOnly) UnmarshalJSON(data []byte) (err error) {
	do.Time, err = time.Parse(time.DateOnly, string(data[1:len(data)-1]))
	if err != nil {
		do.Time, err = time.Parse("02.01.06", string(data[1:len(data)-1]))
	}
	return
}

func (do *DateOnly) ToTime() time.Time {
	return do.Add(0)
}