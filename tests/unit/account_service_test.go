package account_test

import (
	"testing"

	accountv1 "github.com/galadeat/bank-sim/api/proto/account/v1"
	"github.com/galadeat/bank-sim/server/account"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddAccount_Success(t *testing.T) {
	server := account.NewServer()

	ctx := t.Context()

	accID, err := server.AddAccount(ctx, &accountv1.AccountInfo{
		Login: "galadeat",
		Email: "galadeat@xyz.com",
	})

	require.NoError(t, err)
	require.NotEmpty(t, accID.Value)
}

func TestAddAccount_EmptyEmail(t *testing.T) {
	server := account.NewServer()

	ctx := t.Context()

	_, err := server.AddAccount(ctx, &accountv1.AccountInfo{
		Login: "galadeat",
		Email: "",
	})

	require.Error(t, err, "expecting  \"empty email\" error ")

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}

func TestAddAccount_EmptyLogin(t *testing.T) {
	server := account.NewServer()

	ctx := t.Context()

	_, err := server.AddAccount(ctx, &accountv1.AccountInfo{
		Login: "",
		Email: "galadeat@xyz.com",
	})

	require.Error(t, err, "expecting \"empty login\" error")
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}

func TestGetAccount_NotFound(t *testing.T) {
	server := account.NewServer()

	ctx := t.Context()

	_, err := server.GetAccount(ctx, &accountv1.AccountID{Value: "non-existent"})
	require.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
}
