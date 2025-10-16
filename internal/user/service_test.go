package user

import (
	"context"
	"testing"

	userv1 "github.com/galadeat/bank-sim/api/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name        string
		req         *userv1.CreateUserRequest
		wantErr     bool
		wantErrCode codes.Code
	}{
		{
			name: "success",
			req: &userv1.CreateUserRequest{
				Email: "test@test.com",
				Login: "test"},
			wantErr: false,
		},
		{
			name: "miss login",
			req: &userv1.CreateUserRequest{
				Email: "test@test.com",
				Login: "",
			},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "miss email",
			req: &userv1.CreateUserRequest{
				Email: "",
				Login: "test",
			},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			server := NewUserService()

			_, err := server.CreateUser(ctx, tt.req)
			if err != nil {

				if (err != nil) != tt.wantErr {
					t.Errorf("expected error to be %v, got %v", tt.wantErr, err)
				}

				st, _ := status.FromError(err)
				if st.Code() != tt.wantErrCode {
					t.Errorf("expected %v, gt %v", tt.wantErrCode, st.Code())
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	server := NewUserService()

	id, _ := server.CreateUser(ctx, &userv1.CreateUserRequest{
		Login: "test",
		Email: "test@test.com"})

	tests := []struct {
		name        string
		id          string
		wantErr     bool
		wantErrCode codes.Code
	}{
		{
			name:    "success",
			id:      id.Id,
			wantErr: false,
		},
		{
			name:        "user not found",
			id:          "not found",
			wantErr:     true,
			wantErrCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := server.GetUser(ctx, &userv1.GetUserRequest{Id: tt.id})

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error to be %v, got %v", tt.wantErr, err)
			}

			if err != nil {
				st, _ := status.FromError(err)
				if st.Code() != tt.wantErrCode {
					t.Fatalf("expected %v, got %v", tt.wantErrCode, st.Code())
				}
				return
			}

			if user.Login != "test" {
				t.Errorf("expected login \"test\", got %v", user.Login)
			}

			if user.Email != "test@test.com" {
				t.Errorf("expected email \"test@test.com\", got %v", user.Email)
			}
		})
	}

}

func TestListUsers(t *testing.T) {
	t.Run("succes", func(t *testing.T) {
		ctx := context.Background()
		server := NewUserService()

		server.CreateUser(ctx, &userv1.CreateUserRequest{
			Login: "user1",
			Email: "user1@test.com",
		})

		server.CreateUser(ctx, &userv1.CreateUserRequest{
			Login: "user2",
			Email: "user2@test.com",
		})

		res, err := server.ListUsers(ctx, &userv1.ListUsersRequest{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(res.Users) != 2 {
			t.Errorf("expected 2 users, got %v", len(res.Users))
		}

	})

	t.Run("empty list", func(t *testing.T) {
		ctx := context.Background()
		server := NewUserService()

		res, err := server.ListUsers(ctx, &userv1.ListUsersRequest{})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(res.Users) != 0 {
			t.Errorf("expected 0 users, got %v", len(res.Users))
		}

	})

	t.Run("context cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		server := NewUserService()
		_, err := server.ListUsers(ctx, &userv1.ListUsersRequest{})

		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		st, _ := status.FromError(err)

		if st.Code() != codes.Canceled {
			t.Errorf("expected %v, got %v", codes.Canceled, st.Code())
		}
	})

}

func TestUpdateUser(t *testing.T) {
	t.Run("success update both fields", func(t *testing.T) {
		ctx := context.Background()
		server := NewUserService()

		created, _ := server.CreateUser(ctx, &userv1.CreateUserRequest{
			Login: "old_login",
			Email: "old@test.com",
		})

		req := &userv1.UpdateUserRequest{
			Id:    created.Id,
			Email: &wrapperspb.StringValue{Value: "new@test.com"},
			Login: &wrapperspb.StringValue{Value: "new_login"},
		}

		res, err := server.UpdateUser(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.User.Email != "new@test.com" {
			t.Errorf("expected email new@test.com, got %v", res.User.Email)
		}
		if res.User.Login != "new_login" {
			t.Errorf("expected login new_login, got %v", res.User.Login)
		}
	})

	t.Run("update only email", func(t *testing.T) {
		ctx := context.Background()
		server := NewUserService()

		created, _ := server.CreateUser(ctx, &userv1.CreateUserRequest{
			Login: "login",
			Email: "old@test.com",
		})

		req := &userv1.UpdateUserRequest{
			Id:    created.Id,
			Email: &wrapperspb.StringValue{Value: "new@test.com"},
		}

		res, err := server.UpdateUser(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.User.Email != "new@test.com" {
			t.Errorf("expected email new@test.com, got %v", res.User.Email)
		}
		if res.User.Login != "login" {
			t.Errorf("expected login login, got %v", res.User.Login)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		ctx := context.Background()
		server := NewUserService()

		req := &userv1.UpdateUserRequest{
			Id:    "nonexistent",
			Email: &wrapperspb.StringValue{Value: "new@test.com"},
		}

		_, err := server.UpdateUser(ctx, req)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		st, _ := status.FromError(err)
		if st.Code() != codes.NotFound {
			t.Errorf("expected NotFound, got %v", st.Code())
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		server := NewUserService()
		created, _ := server.CreateUser(context.Background(), &userv1.CreateUserRequest{
			Login: "login",
			Email: "email@test.com",
		})

		req := &userv1.UpdateUserRequest{
			Id:    created.Id,
			Email: &wrapperspb.StringValue{Value: "new@test.com"},
		}

		_, err := server.UpdateUser(ctx, req)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		st, _ := status.FromError(err)
		if st.Code() != codes.Canceled {
			t.Errorf("expected Canceled, got %v", st.Code())
		}
	})

	t.Run("empty values should not overwrite", func(t *testing.T) {
		ctx := context.Background()
		server := NewUserService()

		created, _ := server.CreateUser(ctx, &userv1.CreateUserRequest{
			Login: "login",
			Email: "email@test.com",
		})

		req := &userv1.UpdateUserRequest{
			Id:    created.Id,
			Email: &wrapperspb.StringValue{Value: ""},
			Login: &wrapperspb.StringValue{Value: ""},
		}

		res, err := server.UpdateUser(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.User.Email != "email@test.com" {
			t.Errorf("expected email unchanged, got %v", res.User.Email)
		}
		if res.User.Login != "login" {
			t.Errorf("expected login unchanged, got %v", res.User.Login)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("success delete", func(t *testing.T) {
		ctx := context.Background()
		server := NewUserService()

		created, _ := server.CreateUser(ctx, &userv1.CreateUserRequest{
			Login: "login",
			Email: "email@test.com",
		})

		res, err := server.DeleteUser(ctx, &userv1.DeleteUserRequest{Id: created.Id})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !res.Success {
			t.Errorf("expected success=true, got %v", res.Success)
		}

		_, err = server.GetUser(ctx, &userv1.GetUserRequest{Id: created.Id})
		if err == nil {
			t.Errorf("expected error after deletion, got nil")
		}
	})

	t.Run("user not found", func(t *testing.T) {
		ctx := context.Background()
		server := NewUserService()

		res, err := server.DeleteUser(ctx, &userv1.DeleteUserRequest{Id: "nonexistent"})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if res.Success {
			t.Errorf("expected success=false, got %v", res.Success)
		}

		st, _ := status.FromError(err)
		if st.Code() != codes.NotFound {
			t.Errorf("expected NotFound, got %v", st.Code())
		}
	})

	t.Run("context canceled", func(t *testing.T) {

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		server := NewUserService()
		created, _ := server.CreateUser(context.Background(), &userv1.CreateUserRequest{
			Login: "login",
			Email: "email@test.com",
		})

		_, err := server.DeleteUser(ctx, &userv1.DeleteUserRequest{Id: created.Id})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		st, _ := status.FromError(err)
		if st.Code() != codes.Canceled {
			t.Errorf("expected Canceled, got %v", st.Code())
		}
	})
}
