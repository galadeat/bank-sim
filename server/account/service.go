package account

import (
	"context"
	"log"
	"sync"

	accountv2 "github.com/galadeat/bank-sim/api/proto/account/v2"
	commonv1 "github.com/galadeat/bank-sim/api/proto/common/v1"
	userv1 "github.com/galadeat/bank-sim/api/proto/user/v1"
	"github.com/gofrs/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	accountv2.UnimplementedAccountServer
	mu           sync.RWMutex
	accounts     map[string]*accountv2.AccountInfo
	deposits     map[string]*accountv2.DepositResponse
	withdraws    map[string]*accountv2.WithdrawResponse
	accCreations map[string]*accountv2.CreateAccountResponse

	userClient userv1.UserClient
}

// New is the constructor
func New(userClient userv1.UserClient) *Service {
	return &Service{
		accounts:     make(map[string]*accountv2.AccountInfo),
		accCreations: make(map[string]*accountv2.CreateAccountResponse),
		deposits:     make(map[string]*accountv2.DepositResponse),
		withdraws:    make(map[string]*accountv2.WithdrawResponse),
		userClient:   userClient,
	}
}

// CreateAccount is the realization of the rpc method
func (s *Service) CreateAccount(ctx context.Context, req *accountv2.CreateAccountRequest) (*accountv2.CreateAccountResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request canceled by client")
	default:
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	if req.RequestId == "" {
		return nil, status.Error(codes.InvalidArgument, "request id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if resp, ok := s.accCreations[req.RequestId]; ok {
		return resp, nil
	}

	userResp, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: req.UserId})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}
		return nil, status.Errorf(codes.Internal, "failed to call UserService: %v", err)
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error while generating account id: %v", err)
	}

	account := &accountv2.AccountInfo{Id: id.String(),
		Owner:   userResp.User,
		Balance: req.InitialBalance}

	s.accounts[id.String()] = account

	log.Printf("account created: account=%v, request_id=%s", account, req.RequestId)

	resp := &accountv2.CreateAccountResponse{Account: account}
	s.accCreations[req.RequestId] = resp
	return resp, nil
}

// GetAccount is the realization of the rpc method
func (s *Service) GetAccount(ctx context.Context, req *accountv2.GetAccountRequest) (*accountv2.GetAccountResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request canceled by client")
	default:
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "account id is required")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	account, ok := s.accounts[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "account not found")
	}
	return &accountv2.GetAccountResponse{Account: account}, nil

}

// ListAccounts is the realization of the rpc method
func (s *Service) ListAccounts(ctx context.Context, req *accountv2.ListAccountsRequest) (*accountv2.ListAccountsResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request canceled by client")
	default:
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	_, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: req.UserId})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}
		return nil, status.Errorf(codes.Internal, "failed to call UserService: %v", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var accounts []*accountv2.AccountInfo
	for _, account := range s.accounts {
		if account.Owner.Id == req.UserId {
			accounts = append(accounts, account)
		}
	}

	return &accountv2.ListAccountsResponse{Accounts: accounts}, nil
}

// DeleteAccount is the realization of the rpc method
func (s *Service) DeleteAccount(ctx context.Context, req *accountv2.DeleteAccountRequest) (*accountv2.DeleteAccountResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request canceled by client")
	default:
	}

	if req.AccountId == "" {
		return nil, status.Error(codes.InvalidArgument, "account id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	acc, ok := s.accounts[req.AccountId]
	if !ok {
		return nil, status.Error(codes.NotFound, "account not found")
	}

	if acc.Balance.Units != 0 || acc.Balance.Nanos != 0 {
		return nil, status.Error(codes.FailedPrecondition, "cannot delete account with non-zero balance")
	}

	delete(s.accounts, req.AccountId)

	log.Printf("account deleted: id=%s", req.AccountId)

	return &accountv2.DeleteAccountResponse{AccountId: req.AccountId}, nil
}

// Deposit is the realization of the rpc method
func (s *Service) Deposit(ctx context.Context, req *accountv2.DepositRequest) (*accountv2.DepositResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request canceled by client")
	default:
	}
	if req.AccountId == "" {
		return nil, status.Error(codes.InvalidArgument, "account id is required")
	}
	if req.Amount == nil || (req.Amount.Units == 0 && req.Amount.Nanos == 0) {
		return nil, status.Error(codes.InvalidArgument, "deposit must be greater than zero")
	}
	if req.RequestId == "" {
		return nil, status.Error(codes.InvalidArgument, "request id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if resp, ok := s.deposits[req.RequestId]; ok {
		return resp, nil
	}

	acc, ok := s.accounts[req.AccountId]
	if !ok {
		return nil, status.Error(codes.NotFound, "account not found")
	}

	balance, err := addMoney(acc.Balance, req.Amount)
	if err != nil {
		return nil, err
	}

	acc.Balance = balance

	log.Printf("deposit: account_id=%s, request_id=%s, amount=%v, new_balance=%v", req.AccountId, req.RequestId, req.Amount, acc.Balance)

	resp := &accountv2.DepositResponse{Account: acc}
	s.deposits[req.RequestId] = resp
	return resp, nil
}

// Withdraw is the realization of the rpc method.
func (s *Service) Withdraw(ctx context.Context, req *accountv2.WithdrawRequest) (*accountv2.WithdrawResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.Canceled, "request canceled by client")
	default:
	}
	if req.AccountId == "" {
		return nil, status.Error(codes.InvalidArgument, "account id is required")
	}
	if req.RequestId == "" {
		return nil, status.Error(codes.InvalidArgument, "request id is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if resp, ok := s.withdraws[req.RequestId]; ok {
		return resp, nil
	}
	acc, ok := s.accounts[req.AccountId]
	if !ok {
		return nil, status.Error(codes.NotFound, "account not found")
	}
	balance, err := substractMoney(acc.Balance, req.Amount)
	if err != nil {
		return nil, err
	}

	acc.Balance = balance

	log.Printf("withdraw: account_id=%s, request_id=%s, new_balance=%v", req.AccountId, req.RequestId, acc.Balance)
	resp := &accountv2.WithdrawResponse{Account: acc}
	s.withdraws[req.RequestId] = resp
	return resp, nil

}

func addMoney(a, b *commonv1.Money) (*commonv1.Money, error) {
	if m := isZero(a, b); m != nil {
		return m, nil
	}
	if a.Currency != b.Currency {
		return nil, status.Error(codes.InvalidArgument, "currency mismatch")
	}

	balance := &commonv1.Money{
		Currency: a.Currency,
		Nanos:    a.Nanos + b.Nanos,
		Units:    a.Units + b.Units,
	}

	balance = normalizeMoney(balance)

	return balance, nil

}

func substractMoney(a, b *commonv1.Money) (*commonv1.Money, error) {
	if a.Currency != b.Currency {
		return nil, status.Error(codes.FailedPrecondition, "currency mismatch")
	}
	if a.Units < b.Units || (a.Units == b.Units && a.Nanos < b.Nanos) {
		return nil, status.Error(codes.FailedPrecondition, "insufficient balance")
	}
	balance := &commonv1.Money{
		Currency: a.Currency,
		Units:    a.Units - b.Units,
		Nanos:    a.Nanos - b.Nanos,
	}
	balance = normalizeMoney(balance)
	return balance, nil
}

func isZero(a, b *commonv1.Money) *commonv1.Money {
	if a == nil {
		return b
	}

	if b == nil {
		return a
	}

	return nil
}

func normalizeMoney(m *commonv1.Money) *commonv1.Money {
	if m.Nanos >= 1_000_000_000 || m.Nanos <= -1_000_000_000 {
		m.Units += int64(m.Nanos / 1_000_000_000)
		m.Nanos = m.Nanos % 1_000_000_000
	}
	if m.Nanos < 0 {
		m.Units--
		m.Nanos += 1_000_000_000
	}
	return &commonv1.Money{
		Currency: m.Currency,
		Nanos:    m.Nanos,
		Units:    m.Units,
	}
}
