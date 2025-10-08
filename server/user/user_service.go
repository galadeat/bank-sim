package user

import (
	"context"
	"log"
	"sync"

	userv1 "github.com/galadeat/bank-sim/api/proto/user/v1"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	userv1.UnimplementedUserServer
	mu      sync.RWMutex
	userMap map[string]*userv1.UserInfo
}

// Constructor
func NewUserSevice() *UserService {
	return &UserService{userMap: make(map[string]*userv1.UserInfo)}
}

// realizatiion of CreateUser rpc method
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

	s.mu.Lock()
	defer s.mu.Unlock()

	s.userMap[id.String()] = &userv1.UserInfo{Id: id.String(),
		Login: req.Login, Email: req.Email}

	return &userv1.CreateUserResponse{Id: id.String()}, nil
}

// realization of GetUser rpc method
func (s *UserService) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.UserInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, exists := s.userMap[req.Id]
	if exists {
		return value, nil
	}
	return nil, status.Errorf(codes.NotFound, "user does not exist %v", req.Id)

}

// realization of ListUsers rpc method
func (s *UserService) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Errorf(codes.Canceled, "request canceled: %v", ctx.Err())
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*userv1.UserInfo, 0, len(s.userMap))
	for _, u := range s.userMap {
		users = append(users, u)
	}

	return &userv1.ListUsersResponse{
		Users: users,
	}, nil
}

// realization of UpdateUser rpc method
func (s *UserService) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.userMap[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user doesn't exist!")
	}

	if req.Email != nil && req.Email.Value != "" {
		user.Email = req.Email.Value
	}

	if req.Login != nil && req.Login.Value != "" {
		user.Login = req.Login.Value
	}

	return &userv1.UpdateUserResponse{User: user}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.userMap[req.Id]
	if !ok {
		return &userv1.DeleteUserResponse{Success: false}, status.Errorf(codes.NotFound, "user doesn't exist")
	}

	delete(s.userMap, req.Id)

	log.Printf("user deleted: id=%s", req.Id)

	return &userv1.DeleteUserResponse{Success: true}, nil

}
