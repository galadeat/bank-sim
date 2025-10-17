package clients

import (
	accountv2 "github.com/galadeat/bank-sim/api/proto/account/v2"
	userv1 "github.com/galadeat/bank-sim/api/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	accServiceAddr = "localhost:50051"
	usrServiceAddr = "localhost:50052"
)

type Clients struct {
	userConn    *grpc.ClientConn
	accountConn *grpc.ClientConn

	User    userv1.UserClient
	Account accountv2.AccountClient
}

func New() (*Clients, error) {

	userConn, err := grpc.NewClient(
		usrServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	accConn, err := grpc.NewClient(
		accServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Clients{
		userConn:    userConn,
		accountConn: accConn,
		User:        userv1.NewUserClient(userConn),
		Account:     accountv2.NewAccountClient(accConn),
	}, nil
}

func (c *Clients) Close() {
	if c.userConn != nil {
		c.userConn.Close()
	}

	if c.accountConn != nil {
		c.accountConn.Close()
	}
}
