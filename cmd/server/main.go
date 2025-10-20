package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	accountv2 "github.com/galadeat/bank-sim/api/proto/account/v2"
	userv1 "github.com/galadeat/bank-sim/api/proto/user/v1"
	"github.com/galadeat/bank-sim/internal/account"
	"github.com/galadeat/bank-sim/internal/user"
	"github.com/galadeat/bank-sim/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	accountPort = "localhost:50051"
	userPort    = "localhost:50052"
)

func main() {
	file := logger.Init("appServer.log")
	defer file.Close()
	lisUser, err := net.Listen("tcp", userPort)
	if err != nil {
		panic(err)
	}

	grpcUser := grpc.NewServer()
	userSvc := user.New()
	userv1.RegisterUserServer(grpcUser, userSvc)
	go grpcUser.Serve(lisUser)

	connUser, err := grpc.NewClient(userPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	userClient := userv1.NewUserClient(connUser)

	lisAcc, err := net.Listen("tcp", accountPort)
	if err != nil {
		panic(err)
	}
	grpcAcc := grpc.NewServer()
	accSvc := account.New(userClient)
	accountv2.RegisterAccountServer(grpcAcc, accSvc)
	log.Printf("servers started")
	if err := grpcAcc.Serve(lisAcc); err != nil {
		log.Fatalf("account service failed: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("shutting down")
	grpcUser.GracefulStop()
	grpcAcc.GracefulStop()
}
