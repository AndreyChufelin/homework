package storage

import (
	"errors"
	"time"
)

type Event struct {
	ID                        string
	Title                     string
	Date                      time.Time
	EndDate                   time.Time
	Description               string
	UserID                    string
	advanceNotificationPeriod time.Duration
}

var (
	ErrEndDateTooEarly    = errors.New("EndDate is earlier than Date")
	ErrEventAlreadyExists = errors.New("event already exists")
	ErrEventDoesntExist   = errors.New("event doesn't exist")
)
