package service

import (
	"fmt"
	"time"
)

var ErrLastDripNotFound = fmt.Errorf("last drip not found")
var ErrDripAlreadyExist = fmt.Errorf("drip already exist")

// DripService - updating drip service
type DripService interface {
	// GetLastDrip - get last drip time, if not exist ErrLastDripNotFound returns
	GetLastDrip(address string) (time.Time, error)
	// UpdateLastDrip - updating last drip time
	UpdateLastDrip(address string) error
}
