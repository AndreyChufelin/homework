package storage

import (
	"errors"
	"time"

	"github.com/AndreyChufelin/homework/hw12_13_14_15_calendar/internal/validator"
	"github.com/google/uuid"
)

type Event struct {
	ID                        string        `json:"id"`
	Title                     string        `json:"title"`
	Date                      time.Time     `json:"date"`
	EndDate                   time.Time     `json:"end_date"`
	Description               string        `json:"description"`
	UserID                    string        `json:"user_id"`
	AdvanceNotificationPeriod time.Duration `json:"advance_notification_period"`
	Notified                  bool          `json:"-"`
}

var (
	ErrEndDateTooEarly    = errors.New("EndDate is earlier than Date")
	ErrEventAlreadyExists = errors.New("event already exists")
	ErrEventDoesntExist   = errors.New("event doesn't exist")
	ErrNoEventsFound      = errors.New("no events found")
)

func ValidateEvent(validator validator.Validator, event Event) {
	if event.ID != "" {
		_, err := uuid.Parse(event.ID)
		isIDValid := err != nil
		validator.Check(isIDValid, "id", "not valid uuid")
	}

	_, err := uuid.Parse(event.UserID)
	isUserIDValid := err != nil
	validator.Check(isUserIDValid, "user_id", "not valid uuid")

	validator.Check(len(event.Title) > 30, "title", "too long")
	validator.Check(event.Date.After(event.EndDate), "end_date", "too early")
}
