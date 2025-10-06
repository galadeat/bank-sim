package account

import (
	"context"

	userv1 "github.com/galadeat/bank-sim/api/proto/user/v1"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	userv1.UnimplementedUserServer
	userMap map[string]*userv1.UserInfo
}

func NewUserSevice() *UserService {
	return &UserService{userMap: make(map[string]*userv1.UserInfo)}
}

func (s *UserService) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	if req.Login == "" {
		return nil, status.Errorf(codes.InvalidArgument, "login must not be empty")
	}

	if req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email must not be empty")
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while generating user id: %v", err)
	}

	s.userMap[id.String()] = &userv1.UserInfo{Id: id.String(),
		Login: req.Login, Email: req.Email}

	return &userv1.CreateUserResponse{Id: id.String()}, nil
}

func (s *UserService) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.UserInfo, error) {
	value, exists := s.userMap[req.Id]
	if exists {
		return value, nil
	}
	return nil, status.Errorf(codes.NotFound, "user does not exist %v", req.Id)

}
