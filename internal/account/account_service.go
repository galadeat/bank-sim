package account

import (
	"context"

	accountv1 "github.com/galadeat/bank-sim/api/proto/account/v1"

	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	accountv1.UnimplementedAccountServer
	accountMap map[string]*accountv1.AccountInfo
}

// AddAccount realizes Account.addAccount()
func (s *server) AddAccount(ctx context.Context, in *accountv1.AccountInfo) (*accountv1.AccountID, error) {
	if in.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login cannot be empty")
	}

	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email cannot be empty")
	}

	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error while generating Account ID: %v", err)
	}
	in.Id = out.String()

	s.accountMap[in.Id] = in
	return &accountv1.AccountID{Value: in.Id}, nil
}

// GetAccount realizes  Account.getAccount()
func (s *server) GetAccount(ctx context.Context, in *accountv1.AccountID) (*accountv1.AccountInfo, error) {
	value, exists := s.accountMap[in.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Product does not exist %v", in.Value)
}

// create server when initializing
func NewServer() *server {
	return &server{
		accountMap: make(map[string]*accountv1.AccountInfo),
	}
}
