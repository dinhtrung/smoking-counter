package services

import (
	"github.com/dinhtrung/smoking-counter/internal/app/smoke-counter/web/dto"
)

type SmokeService interface {
	// GetAll Return a sorted list of all Smoke records for a given user
	GetAll(user string) []*dto.DailySmokingDTO

	// Create a new Smoke record for current user and hour of today
	Create(user, hour string) (*dto.DailySmokingDTO, error)

	// Delete an existing Smoke record for a given user based on date and hour
	Delete(user, hour string) error

	// DeleteAll existing Smoke record for a given user
	DeleteAll(user string) error

	// DeleteOne an existing Smoke record for a given user based on date and hour
	DeleteOne(user, date, hour string) error

	// DeleteAllByDate exists Smoke record for a single day
	DeleteAllByDate(user, date string) error

	// DeleteAllBefore delete all Smoke records before a given date
	DeleteAllBefore(user, date string) (int, error)

	// DeleteAllAfter delete all Smoke records after a given date
	DeleteAllAfter(user, date string) (int, error)
}
