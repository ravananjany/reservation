package utils

import (
	"fmt"

	"github.com/reservation/constants"
	"github.com/reservation/protos"
	"github.com/reservation/resources"
)

func ReservationModelMapper(res *protos.Reservation) *resources.Reservation {

	reservation := &resources.Reservation{
		UserId: res.User.UserId,
		From:   constants.BOARD,
		To:     constants.DEST,
		Price:  constants.PRICE,
		Seat:   fmt.Sprintf(res.Booking.GetSection()+" %d", res.Booking.Seat),
		User: resources.User{
			FirstName: res.User.FirstName,
			LastName:  res.User.LastName,
			EmailId:   res.User.EmailId,
		},
	}
	return reservation
}
