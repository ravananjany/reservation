package service

import (
	"context"
	"testing"

	"github.com/reservation/constants"
	"github.com/reservation/protos"
	"github.com/reservation/resources"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *mockReservationClient) CreateTicket(ctx context.Context, in *protos.User, opts ...grpc.CallOption) (*protos.Reservation, error) {
	if r.skip {
		return &protos.Reservation{User: &protos.User{}, Booking: &protos.BookingDetails{}}, nil
	}

	if r.failure {
		return nil, constants.ErrInternalServer
	}
	return &protos.Reservation{User: &protos.User{}, Booking: &protos.BookingDetails{}}, nil
}

func (r mockReservationClient) ViewTicket(ctx context.Context, in *protos.UserId, opts ...grpc.CallOption) (*protos.Reservation, error) {
	if r.notfound {
		return nil, status.Error(codes.NotFound, constants.ErrNotFound.Error())
	}
	return &protos.Reservation{User: &protos.User{}, Booking: &protos.BookingDetails{}}, nil
}

func (r mockReservationClient) Viewreservations(ctx context.Context, in *protos.ReadAll, opts ...grpc.CallOption) (*protos.Reservations, error) {
	if r.failure {
		return nil, constants.ErrNotFound
	}
	return &protos.Reservations{Reservation: []*protos.Reservation{}}, nil
}

func (r *mockReservationClient) DeleteTicket(ctx context.Context, in *protos.UserId, opts ...grpc.CallOption) (*protos.DeleteResponse, error) {
	if r.notfound {
		return nil, status.Error(codes.NotFound, constants.ErrNotFound.Error())
	}
	return &protos.DeleteResponse{MessageResponse: constants.DELETERES}, nil
}

func (r *mockReservationClient) UpdateTicket(ctx context.Context, in *protos.User, opts ...grpc.CallOption) (*protos.Reservation, error) {
	if r.notfound {
		return nil, status.Error(codes.NotFound, constants.ErrNotFound.Error())
	}
	return &protos.Reservation{User: &protos.User{}, Booking: &protos.BookingDetails{}}, nil
}

type mockReservationClient struct {
	protos.ReservationServiceClient
	failure  bool
	notfound bool
	skip     bool
}

func TestCreate(t *testing.T) {
	type feilds struct {
		ctx    context.Context
		user   resources.User
		client *mockReservationClient
	}

	tests := []struct {
		name string
		args feilds
	}{
		{name: "create success", args: feilds{context.Background(), resources.User{}, &mockReservationClient{}}},
		{name: "create failure", args: feilds{context.Background(), resources.User{}, &mockReservationClient{failure: true}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := NewReservationService(tt.args.client, &logrus.Logger{})
			_, err := obj.CreateTicket(tt.args.ctx, tt.args.user)

			if tt.name == "create failure" {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestView(t *testing.T) {
	type feilds struct {
		ctx    context.Context
		user   resources.User
		client *mockReservationClient
	}

	tests := []struct {
		name string
		args feilds
	}{
		{name: "view success", args: feilds{context.Background(), resources.User{}, &mockReservationClient{}}},
		{name: "view failure", args: feilds{context.Background(), resources.User{}, &mockReservationClient{notfound: true}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := NewReservationService(tt.args.client, &logrus.Logger{})
			res, _ := obj.CreateTicket(tt.args.ctx, tt.args.user)
			testres, err := obj.ViewTicket(tt.args.ctx, res.UserId)
			if tt.name == "view failure" {
				assert.Equal(t, err, constants.ErrNotFound)
			} else {
				assert.NotNil(t, testres.UserId)
			}
		})
	}
}
func TestViews(t *testing.T) {
	type feilds struct {
		ctx    context.Context
		user   resources.User
		client *mockReservationClient
	}

	tests := []struct {
		name string
		args feilds
	}{
		{name: "view success", args: feilds{context.Background(), resources.User{}, &mockReservationClient{}}},
		{name: "view failure", args: feilds{context.Background(), resources.User{}, &mockReservationClient{notfound: true}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := NewReservationService(tt.args.client, &logrus.Logger{})
			_, _ = obj.CreateTicket(tt.args.ctx, tt.args.user)
			_, err := obj.ViewReservations(tt.args.ctx, constants.SEC_B)
			assert.Nil(t, err)
		})
	}
}

func TestDelete(t *testing.T) {
	type feilds struct {
		ctx    context.Context
		user   resources.User
		client *mockReservationClient
	}

	tests := []struct {
		name string
		args feilds
	}{
		{name: "Delete success", args: feilds{context.Background(), resources.User{}, &mockReservationClient{}}},
		{name: "Delete failure", args: feilds{context.Background(), resources.User{}, &mockReservationClient{notfound: true}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := NewReservationService(tt.args.client, &logrus.Logger{})
			res, _ := obj.CreateTicket(tt.args.ctx, tt.args.user)
			_, err := obj.DeleteTicket(tt.args.ctx, res.UserId)
			if tt.name == "Delete failure" {
				assert.Equal(t, err, constants.ErrNotFound)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type feilds struct {
		ctx    context.Context
		user   resources.User
		client *mockReservationClient
	}

	tests := []struct {
		name string
		args feilds
	}{
		{name: "Delete success", args: feilds{context.Background(), resources.User{}, &mockReservationClient{}}},
		{name: "Delete failure", args: feilds{context.Background(), resources.User{}, &mockReservationClient{notfound: true}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := NewReservationService(tt.args.client, &logrus.Logger{})
			res, _ := obj.CreateTicket(tt.args.ctx, tt.args.user)
			user := &resources.User{EmailId: "abc@gmail.com"}
			user.UserId = res.UserId
			_, err := obj.UpdateTicket(tt.args.ctx, *user)
			if tt.name == "Delete failure" {
				assert.Equal(t, err, constants.ErrNotFound)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
