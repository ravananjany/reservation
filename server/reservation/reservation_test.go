package reservation

import (
	"context"
	"testing"

	"github.com/reservation/protos"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var testEmailId = "abc@gmail.com"

func GetUser() *protos.User {
	user := &protos.User{
		FirstName: "jany", LastName: "ram", EmailId: testEmailId,
	}
	return user
}

func TestCreateTicket(t *testing.T) {

	obj := NewReservationService(logrus.New())
	res, _ := obj.CreateTicket(context.Background(), GetUser())
	resview, _ := obj.ViewTicket(context.Background(), &protos.UserId{UserId: res.User.UserId})
	assert.Equal(t, testEmailId, resview.User.EmailId)
}

func TestAllReservations(t *testing.T) {
	obj := NewReservationService(logrus.New())
	_, _ = obj.CreateTicket(context.Background(), GetUser())
	res, _ := obj.Viewreservations(context.Background(), &protos.ReadAll{})
	assert.Greater(t, len(res.Reservation), 0)
}

func TestDeleteTicket(t *testing.T) {
	obj := NewReservationService(logrus.New())
	res, _ := obj.CreateTicket(context.Background(), GetUser())
	obj.DeleteTicket(context.Background(), &protos.UserId{UserId: res.User.GetUserId()})
	_, err := obj.ViewTicket(context.Background(), &protos.UserId{UserId: res.User.GetUserId()})
	assert.NotNil(t, err)
}

func TestUpdateTicket(t *testing.T) {
	obj := NewReservationService(logrus.New())
	res, _ := obj.CreateTicket(context.Background(), GetUser())
	user := GetUser()
	user.UserId = res.User.UserId
	user.EmailId = "jan@gmail.com"
	updateres, _ := obj.UpdateTicket(context.Background(), user)
	assert.Equal(t, user.EmailId, updateres.User.GetEmailId())

}
