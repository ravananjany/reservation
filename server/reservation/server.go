package reservation

import (
	"fmt"
	"net"

	"github.com/reservation/config"
	"github.com/reservation/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Startserver() {
	log, conf := config.LoadConfig()
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", conf.Grpcport))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ser := grpc.NewServer()

	protos.RegisterReservationServiceServer(ser, NewReservationService(log))
	reflection.Register(ser)
	log.Info("starting  grpc server")
	if err := ser.Serve(listener); err != nil {
		log.Error(err)
		log.Errorf("unable to start the server %s", conf.Grpcport)
	}
}
