package account

import (
	"context"
	"testing"

	accountv2 "github.com/galadeat/bank-sim/api/proto/account/v2"
	commonv1 "github.com/galadeat/bank-sim/api/proto/common/v1"
	userv1 "github.com/galadeat/bank-sim/api/proto/user/v1"
	"github.com/galadeat/bank-sim/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateAccount(t *testing.T) {
	tests := []struct {
		name        string
		req         *accountv2.CreateAccountRequest
		mockSetup   func(user *mocks.MockUserClient)
		wantErr     bool
		wantErrCode codes.Code
	}{
		{
			name: "success",
			req: &accountv2.CreateAccountRequest{
				UserId: "user-123",
				InitialBalance: &commonv1.Money{
					Currency: "USD",
					Units:    1000,
					Nanos:    0,
				},
				RequestId: "1",
			},
			mockSetup: func(user *mocks.MockUserClient) {
				user.EXPECT().
					GetUser(gomock.Any(), &userv1.GetUserRequest{Id: "user-123"}).
					Return(&userv1.GetUserResponse{User: &userv1.UserInfo{Id: "user-123"}}, nil)
			},
			wantErr:     false,
			wantErrCode: codes.OK,
		},
		{
			name: "userId is empty",
			req: &accountv2.CreateAccountRequest{
				UserId: "",
				InitialBalance: &commonv1.Money{
					Currency: "USD",
					Units:    1000,
					Nanos:    0,
				},
				RequestId: "1",
			},
			mockSetup:   nil,
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "requestId is empty",
			req: &accountv2.CreateAccountRequest{
				UserId: "user-123",
				InitialBalance: &commonv1.Money{
					Currency: "USD",
					Units:    1000,
					Nanos:    0,
				},
				RequestId: "",
			},
			mockSetup:   nil,
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "userId is invalid",
			req: &accountv2.CreateAccountRequest{
				UserId: "user-123",
				InitialBalance: &commonv1.Money{
					Currency: "USD",
					Units:    1000,
					Nanos:    0,
				},
				RequestId: "1",
			},
			mockSetup: func(user *mocks.MockUserClient) {
				user.EXPECT().
					GetUser(gomock.Any(), &userv1.GetUserRequest{Id: "user-123"}).
					Return(nil, status.Error(codes.InvalidArgument, "invalid user id"))
			},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockUserClient(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			svc := New(mock)

			_, err := svc.CreateAccount(context.Background(), tt.req)

			if tt.wantErr && err == nil || !tt.wantErr && err != nil {
				t.Errorf("expected %v, got %v", tt.wantErrCode, err)
			}

		})
	}

	t.Run("repeated request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := mocks.NewMockUserClient(ctrl)
		user.EXPECT().
			GetUser(gomock.Any(), &userv1.GetUserRequest{Id: "user-123"}).
			Return(&userv1.GetUserResponse{User: &userv1.UserInfo{Id: "user-123"}}, nil)

		svc := New(user)

		req := &accountv2.CreateAccountRequest{
			UserId: "user-123",
			InitialBalance: &commonv1.Money{
				Currency: "USD",
				Units:    1000,
				Nanos:    0,
			},
			RequestId: "1",
		}

		firstResp, err := svc.CreateAccount(context.Background(), req)
		assert.Nil(t, err)
		secondResp, err := svc.CreateAccount(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, firstResp, secondResp)
	})

	t.Run("cancelled request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		defer ctrl.Finish()
		user := mocks.NewMockUserClient(ctrl)

		svc := New(user)

		req := &accountv2.CreateAccountRequest{
			UserId: "user-123",
			InitialBalance: &commonv1.Money{
				Currency: "USD",
				Units:    1000,
				Nanos:    0,
			},
			RequestId: "1",
		}
		_, err := svc.CreateAccount(ctx, req)
		assert.Error(t, err)

	})

}

func TestGetAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := mocks.NewMockUserClient(ctrl)
	user.EXPECT().
		GetUser(gomock.Any(), &userv1.GetUserRequest{Id: "user-123"}).
		Return(&userv1.GetUserResponse{User: &userv1.UserInfo{Id: "user-123"}}, nil)

	svc := New(user)

	accCreated, err := svc.CreateAccount(context.Background(), &accountv2.CreateAccountRequest{
		UserId: "user-123",
		InitialBalance: &commonv1.Money{
			Currency: "USD",
			Units:    1000,
			Nanos:    0,
		},
		RequestId: "1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name        string
		req         *accountv2.GetAccountRequest
		wantErr     bool
		wantErrCode codes.Code
	}{
		{
			name:        "success",
			req:         &accountv2.GetAccountRequest{Id: accCreated.Account.GetId()},
			wantErr:     false,
			wantErrCode: codes.OK,
		},
		{
			name:        "accountId is empty",
			req:         &accountv2.GetAccountRequest{Id: ""},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name:        "accountId is missing",
			req:         &accountv2.GetAccountRequest{Id: "not-found"},
			wantErr:     true,
			wantErrCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accGot, err := svc.GetAccount(context.Background(), tt.req)

			if tt.wantErr && err == nil || !tt.wantErr && err != nil {
				t.Errorf("expected %v, got %v", tt.wantErrCode, err)
			}
			if accGot != nil {
				assert.Equal(t, accGot.Account, accCreated.Account)
			}

		})
	}
}
