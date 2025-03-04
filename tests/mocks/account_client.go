package mocks

import (
	"context"
	"flame/pkg/pb"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockAccountClient struct {
	mock.Mock
}

func (mock *MockAccountClient) Register(ctx context.Context, r *pb.RegisterReq, opts ...grpc.CallOption) (*pb.RegisterRes, error) {
	args := mock.Called(ctx, r, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.RegisterRes), args.Error(1)
}
func (mock *MockAccountClient) Login(ctx context.Context, r *pb.LoginReq, opts ...grpc.CallOption) (*pb.LoginRes, error) {
	args := mock.Called(ctx, r, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.LoginRes), args.Error(1)
}
func (mock *MockAccountClient) GetTokens(ctx context.Context, r *pb.GetTokensReq, opts ...grpc.CallOption) (*pb.GetTokensRes, error) {
	args := mock.Called(ctx, r, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GetTokensRes), args.Error(1)
}
