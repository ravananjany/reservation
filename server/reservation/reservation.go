package reservation

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/reservation/constants"
	"github.com/reservation/protos"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReservationService struct {
	protos.ReservationServiceServer
	reservationMap map[string]*protos.Reservation
	sectionA       int
	sectionB       int
	mutex          sync.Mutex
	logger         *logrus.Logger
}

func NewReservationService(log *logrus.Logger) *ReservationService {
	return &ReservationService{
		reservationMap: make(map[string]*protos.Reservation),
		logger:         log,
	}
}

func (r *ReservationService) CreateTicket(ctx context.Context, user *protos.User) (*protos.Reservation, error) {
	r.logger.Info("enter CreateTicket server")

	r.mutex.Lock()
	defer r.mutex.Unlock()

	userId := uuid.New().String()
	user.UserId = userId
	booking := protos.BookingDetails{}

	if r.sectionA < r.sectionB {
		r.sectionA++
		booking.Section = constants.SEC_A
		booking.Seat = int32(r.sectionA + 1)
	} else {
		r.sectionB++
		booking.Section = constants.SEC_B
		booking.Seat = int32(r.sectionB + 1)
	}

	dp, err := r.CheckPriceDiscount(user.DiscountCode)
	if err != nil {
		return nil, err
	}

	booking.DiscountPrice = dp

	reservation := &protos.Reservation{
		Booking: &booking,
		User:    user,
	}
	r.reservationMap[user.UserId] = reservation

	return reservation, nil
}

func (r *ReservationService) CheckPriceDiscount(dis string) (int32, error) {
	if len(dis) == 0 {
		return 20, nil
	}

	//discount30
	//dthirty

	var sb strings.Builder
	for i, v := range dis {
		if i != 0 {
			sb.WriteRune(v)
		}
	}

	dp, err := strconv.Atoi(sb.String())
	if err != nil {
		r.logger.Error(err)
		return 0, status.Error(codes.InvalidArgument, constants.ErrInvalidDiscount.Error())
	}
	if dp > 20 {
		return 0, status.Error(codes.OutOfRange, constants.ErrDiscountOutRange.Error())
	}

	fd := 20 - dp

	return int32(fd), nil
}

func (r *ReservationService) ViewTicket(ctx context.Context, id *protos.UserId) (*protos.Reservation, error) {
	r.logger.Info("enter ViewTicket server")
	var res *protos.Reservation
	var found bool
	r.mutex.Lock()
	defer r.mutex.Unlock()

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
	r.logger.Info("exit")
	return nil, status.Error(codes.NotFound, constants.ErrNotFound.Error())
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
	defer r.mutex.Unlock()
	var found bool
	for K := range r.reservationMap {
		if userId.UserId == K {
			found = true
			break
		}
	}
	if !found {
		return nil, status.Error(codes.NotFound, constants.ErrNotFound.Error())
	}
	delete(r.reservationMap, userId.UserId)

	return &protos.DeleteResponse{MessageResponse: constants.DELETERES}, nil
}

func (r *ReservationService) UpdateTicket(ctx context.Context, user *protos.User) (*protos.Reservation, error) {
	r.logger.Info("enter UpdateTicket server")

	r.mutex.Lock()
	defer r.mutex.Unlock()
	currentUserResDetails := GetUserDetails(r.reservationMap, user.UserId)
	if currentUserResDetails == nil {
		return nil, status.Error(codes.NotFound, constants.ErrNotFound.Error())
	}
	delete(r.reservationMap, currentUserResDetails.User.UserId)
	currentUserResDetails.User = user
	r.reservationMap[user.UserId] = currentUserResDetails

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
