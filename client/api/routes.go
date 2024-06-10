package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/reservation/client/service"
	"github.com/reservation/config"
	"github.com/reservation/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartRest(port string) {
	router := gin.Default()
	router.POST("/ticket", CreateTicket)
	router.GET("/ticket/:id", ViewTicket)
	router.Run(fmt.Sprintf("localhost:%s", port))
}
func StartClient() {
	log, conf := config.LoadConfig()
	logger = log
	log.Info("connecting client")
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", conf.Grpcport), grpc.WithTransportCredentials(insecure.NewCredentials()))
	//conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", conf.Grpcport), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error connecting to server", err)
	}
	defer conn.Close()
	reservationClient = service.NewReservationService(protos.NewReservationServiceClient(conn), log)
	log.Info("starting rest api")
	StartRest(conf.Apiport)
}
