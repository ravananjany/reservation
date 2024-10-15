package constants

import "errors"

var (
	ErrNotFound         = errors.New("ticket not found")
	ErrPubNotPermitted  = errors.New("invalid in input data publication updation not permitted")
	ErrReservationFull  = errors.New("no seats available")
	ErrInternalServer   = errors.New("internal server error")
	ErrDiscountOutRange = errors.New("discount price is too high")
	ErrInvalidDiscount  = errors.New("invalid discount")
)
