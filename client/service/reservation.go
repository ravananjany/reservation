package service

import (
	"context"
	"strings"

	"github.com/reservation/constants"
	"github.com/reservation/protos"
	"github.com/reservation/resources"
	"github.com/reservation/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReservationInterface interface {
	CreateTicket(ctx context.Context, user resources.User) (*resources.Reservation, error)
	ViewTicket(ctx context.Context, userId string) (*resources.Reservation, error)
	ViewReservations(ctx context.Context, section string) ([]*resources.Reservation, error)
	DeleteTicket(ctx context.Context, userId string) (string, error)
	UpdateTicket(ctx context.Context, user resources.User) (*resources.Reservation, error)
}

type reservation struct {
	reservationServer protos.ReservationServiceClient
	logger            *logrus.Logger
}

func NewReservationService(rs protos.ReservationServiceClient, log *logrus.Logger) ReservationInterface {
	return &reservation{logger: log, reservationServer: rs}
}

func (r *reservation) CreateTicket(ctx context.Context, user resources.User) (*resources.Reservation, error) {
	r.logger.Info("enter Create ticket client")
	protoUser := &protos.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		EmailId:   user.EmailId,
	}
	res, err := r.reservationServer.CreateTicket(ctx, protoUser)
	if err != nil {
		return nil, err
	}
	reservation := utils.ReservationModelMapper(res)
	return reservation, nil
}

func (r *reservation) ViewTicket(ctx context.Context, userId string) (*resources.Reservation, error) {
	r.logger.Info("enter ViewTicket  client")
	protoUserId := &protos.UserId{UserId: userId}
	res, err := r.reservationServer.ViewTicket(ctx, protoUserId)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, constants.ErrNotFound
		}
	}

	return utils.ReservationModelMapper(res), nil
}

func (r *reservation) ViewReservations(ctx context.Context, section string) ([]*resources.Reservation, error) {
	r.logger.Info("enter ViewReservations client")
	readAll := &protos.ReadAll{Section: strings.ToUpper(section)}
	res, err := r.reservationServer.Viewreservations(ctx, readAll)
	if err != nil {
		return nil, constants.ErrNotFound
	}
	result := []*resources.Reservation{}
	for _, v := range res.Reservation {
		result = append(result, utils.ReservationModelMapper(v))
	}
	return result, nil
}

func (r *reservation) DeleteTicket(ctx context.Context, userId string) (string, error) {
	r.logger.Info("enter DeleteTicket client")
	protoUserId := &protos.UserId{UserId: userId}
	res, err := r.reservationServer.DeleteTicket(ctx, protoUserId)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return "nil", constants.ErrNotFound
		}
	}
	return res.MessageResponse, nil
}

func (r *reservation) UpdateTicket(ctx context.Context, user resources.User) (*resources.Reservation, error) {
	r.logger.Info("enter DeleteTicket client")
	protoUser := &protos.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		EmailId:   user.EmailId,
		UserId:    user.UserId,
	}
	res, err := r.reservationServer.UpdateTicket(ctx, protoUser)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, constants.ErrNotFound
		}
	}
	return utils.ReservationModelMapper(res), nil
}
