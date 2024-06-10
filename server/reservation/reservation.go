package reservation

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/reservation/constants"
	"github.com/reservation/protos"
	"github.com/sirupsen/logrus"
)

type ReservationService struct {
	protos.ReservationServiceServer
	reservationMap map[string]*protos.Reservation
	bookingMap     map[string]int //map that holds the seat number
	sectionA       int
	sectionB       int
	mutex          sync.Mutex
	logger         *logrus.Logger
}

func NewReservationService(log *logrus.Logger) *ReservationService {
	return &ReservationService{
		reservationMap: make(map[string]*protos.Reservation),
		bookingMap:     make(map[string]int),
		logger:         log,
	}
}

func (r *ReservationService) CreateTicket(ctx context.Context, user *protos.User) (*protos.Reservation, error) {
	r.logger.Info("enter CreateTicket server")

	r.mutex.Lock()
	userId := uuid.New().String()
	user.UserId = userId
	booking := protos.BookingDetails{}

	if r.sectionA+r.sectionB >= constants.SEAT_LIMIT {
		return nil, constants.ErrReservationFull
	}

	if r.sectionA < r.sectionB {
		r.sectionA++
		booking.Section = constants.SEC_A
		booking.Seat = int32(r.sectionA + 1)
	} else {
		r.sectionB++
		booking.Section = constants.SEC_B
		booking.Seat = int32(r.sectionB + 1)
	}

	reservation := &protos.Reservation{
		Booking: &booking,
		User:    user,
	}
	r.reservationMap[user.UserId] = reservation

	r.mutex.Unlock()
	return reservation, nil
}

func (r *ReservationService) ViewTicket(ctx context.Context, id *protos.UserId) (*protos.Reservation, error) {
	r.logger.Info("enter ViewTicket server")
	var res *protos.Reservation
	var found bool
	r.mutex.Lock()

	for k, V := range r.reservationMap {
		if k == id.UserId {
			res = V
			found = true
			break
		}
	}

	if found {
		return res, nil
	}
	r.mutex.Unlock()

	return nil, constants.ErrNotFound
}

func (r *ReservationService) Viewreservations(ctx context.Context, all *protos.ReadAll) (*protos.Reservations, error) {
	r.logger.Info("enter Viewreservations server")

	res := []*protos.Reservation{}
	resB := []*protos.Reservation{}
	for _, v := range r.reservationMap {
		switch v.Booking.GetSection() {
		case constants.SEC_A:
			res = append(res, v)
		case constants.SEC_B:
			resB = append(resB, v)

		}
	}

	if all.Section == constants.SEC_A {
		return &protos.Reservations{Reservation: res}, nil
	} else if all.Section == constants.SEC_B {
		return &protos.Reservations{Reservation: resB}, nil
	}

	res = append(res, resB...)
	return &protos.Reservations{Reservation: res}, nil

}

func (r *ReservationService) DeleteTicket(ctx context.Context, userId *protos.UserId) (*protos.DeleteResponse, error) {
	r.logger.Info("enter DeleteTicket server")

	r.mutex.Lock()
	var found bool
	for K := range r.reservationMap {
		if userId.UserId == K {
			found = true
			break
		}
	}
	if !found {
		return nil, constants.ErrNotFound
	}

	delete(r.reservationMap, userId.UserId)

	r.mutex.Unlock()

	return &protos.DeleteResponse{MessageResponse: constants.DELETERES}, nil
}

func (r *ReservationService) UpdateTicket(ctx context.Context, user *protos.User) (*protos.Reservation, error) {
	r.logger.Info("enter UpdateTicket server")

	r.mutex.Lock()
	currentUserResDetails := GetUserDetails(r.reservationMap, user.UserId)
	if currentUserResDetails == nil {
		return nil, constants.ErrNotFound
	}
	delete(r.reservationMap, currentUserResDetails.User.UserId)
	currentUserResDetails.User = user
	r.reservationMap[user.UserId] = currentUserResDetails
	r.mutex.Unlock()
	return currentUserResDetails, nil
}

func GetUserDetails(mp map[string]*protos.Reservation, userId string) *protos.Reservation {

	for K, V := range mp {
		if K == userId {
			return V
		}
	}
	return nil
}
