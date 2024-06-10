package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reservation/client/service"
	"github.com/reservation/constants"
	"github.com/reservation/resources"
	"github.com/sirupsen/logrus"
)

var reservationClient service.ReservationInterface
var logger *logrus.Logger

func CreateTicket(c *gin.Context) {
	funcDesc := "CreateTicket"
	logger.Info("enter  " + funcDesc)
	var user resources.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := reservationClient.CreateTicket(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func ViewTicket(c *gin.Context) {
	funcDesc := "ViewTicket"
	logger.Info("enter  " + funcDesc)
	UserId := c.Param("id")
	res, err := reservationClient.ViewTicket(c.Request.Context(), UserId)
	if err.Error() == constants.ErrNotFound.Error() {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func DeleteTicket(c *gin.Context) {
	funcDesc := "DeleteTicket"
	logger.Info("enter  " + funcDesc)
	UserId := c.Param("id")
	res, err := reservationClient.DeleteTicket(c.Request.Context(), UserId)
	if err == constants.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, res)
}

func ViewAllReservations(c *gin.Context) {
	funcDesc := "ViewAllReservations"
	logger.Info("enter  " + funcDesc)
	section := c.Param("section")
	res, err := reservationClient.ViewReservations(c.Request.Context(), section)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func UpdateTicket(c *gin.Context) {
	funcDesc := "UpdateTicket"
	logger.Info("enter  " + funcDesc)
	var user resources.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := reservationClient.UpdateTicket(c.Request.Context(), user)
	if err == constants.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, res)
}
