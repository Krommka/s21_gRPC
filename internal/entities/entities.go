package entities

import "time"

type Entry struct {
	SessionId string
	Frequency float64
	Timestamp time.Time
}
