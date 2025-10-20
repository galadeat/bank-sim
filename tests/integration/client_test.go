package client_test

// TODO: need to rewrite tests

//import (
//	"context"
//	"net"
//	"testing"
//	"time"
//
//	accountv1 "github.com/galadeat/bank-sim/api/proto/account/v1"
//	"github.com/galadeat/bank-sim/server/account"
//	"github.com/stretchr/testify/require"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials/insecure"
//)
//
//func startTestServer(t *testing.T) (accountv1.AccountClient, func()) {
//	lis, err := net.Listen("tcp", ":0")
//	require.NoError(t, err)
//
//	grpcServer := grpc.NewServer()
//	accountv1.RegisterAccountServer(grpcServer, account.NewServer())
//
//	go grpcServer.Serve(lis)
//
//	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
//
//	require.NoError(t, err)
//
//	client := accountv1.NewAccountClient(conn)
//
//	cleanup := func() {
//		grpcServer.Stop()
//		conn.Close()
//		lis.Close()
//	}
//
//	return client, cleanup
//}
//
//func TestClient_CreateAndGetAccount(t *testing.T) {
//	client, cleanup := startTestServer(t)
//	defer cleanup()
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//
//	createResp, err := client.AddAccount(ctx, &accountv1.AccountInfo{
//		Login: "galadeat",
//		Email: "galadeat@xyz.com",
//	})
//	require.NoError(t, err)
//	require.NotEmpty(t, createResp.Value)
//
//	getResp, err := client.GetAccount(ctx, &accountv1.AccountID{
//		Value: createResp.Value,
//	})
//	require.NoError(t, err)
//	require.Equal(t, "galadeat", getResp.Login)
//	require.Equal(t, "galadeat@xyz.com", getResp.Email)
//
//}
