package constants

import "errors"

var (
	ErrNotFound        = errors.New("ticket not found")
	ErrPubNotPermitted = errors.New("invalid in input data publication updation not permitted")
	ErrReservationFull = errors.New("no seats available")
)
