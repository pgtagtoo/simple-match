package models

import (
	"github.com/google/uuid"
	"time"
)

type Person struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Height      int       `json:"height"`
	Gender      string    `json:"gender"`
	WantedDates int       `json:"wanted_dates"`
	CreatedAt   time.Time
}
