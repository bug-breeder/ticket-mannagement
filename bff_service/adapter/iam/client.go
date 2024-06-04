package iam

import (
	"context"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"

	"go.tekoapis.com/tekone/app/warehouse/iam_service/api"
)

type ClientAdapter interface {
	CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error)
	GetToken(ctx context.Context, req *api.GetTokenRequest) (*api.GetTokenResponse, error)
	GetUserInfo(ctx context.Context, req *api.GetUserInfoRequest) (*api.GetUserInfoResponse, error)
}

func New(log logr.Logger, addr string) (ClientAdapter, error) {
	tmConn, err := grpc.DialContext(context.Background(), addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	iamClient := api.NewIamServiceClient(tmConn)
	return &client{
		client: iamClient,
		log:    log,
	}, nil
}

type client struct {
	client api.IamServiceClient
	log    logr.Logger
}

func (c *client) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	return c.client.CreateUser(ctx, req)
}

func (c *client) GetToken(ctx context.Context, req *api.GetTokenRequest) (*api.GetTokenResponse, error) {
	return c.client.GetToken(ctx, req)
}

func (c *client) GetUserInfo(ctx context.Context, req *api.GetUserInfoRequest) (*api.GetUserInfoResponse, error) {
	return c.client.GetUserInfo(ctx, req)
}
