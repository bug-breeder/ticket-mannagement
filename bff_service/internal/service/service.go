package service

import (
	"context"

	"go.tekoapis.com/tekone/app/warehouse/bff_service/adapter/iam"
	"go.tekoapis.com/tekone/app/warehouse/bff_service/adapter/tm"
	"go.tekoapis.com/tekone/app/warehouse/bff_service/api"
	"go.tekoapis.com/tekone/app/warehouse/bff_service/internal/store"
	iam_api "go.tekoapis.com/tekone/app/warehouse/iam_service/api"
	tm_api "go.tekoapis.com/tekone/app/warehouse/tm_service/api"
	health "go.tekoapis.com/tekone/library/grpc/health"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.tekoapis.com/tekone/app/warehouse/bff_service/config"

	"github.com/go-logr/logr"
)

type Service struct {
	log logr.Logger
	// more connector here
	store store.StoreQuerier
	// embedded unimplemented service server
	health.UnimplementedHealthCheckServiceServer
	api.UnimplementedBffServiceServer
	iamClient iam.ClientAdapter
	tmClient  tm.ClientAdapter
}

func NewService(
	logger logr.Logger,
	cfg *config.Config,
) (*Service, error) {

	iamClient, err := iam.New(logger, cfg.IAMHost)
	if err != nil {
		return nil, err
	}

	tmClient, err := tm.New(logger, cfg.TMHost)
	if err != nil {
		return nil, err
	}

	return &Service{
		log:       logger,
		iamClient: iamClient,
		tmClient:  tmClient,
	}, nil
}

func (s *Service) CreateUser(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	iamReq := &iam_api.CreateUserRequest{
		FullName:  req.FullName,
		Username:  req.Username,
		Gender:    req.Gender,
		BirthDate: req.BirthDate,
		Password:  req.Password,
		Role:      req.Role,
	}

	resp, err := s.iamClient.CreateUser(ctx, iamReq)
	if err != nil {
		return nil, err
	}

	return &api.CreateUserResponse{
		Message: resp.Message,
		UserId:  resp.UserId,
	}, nil
}

func (s *Service) GetToken(ctx context.Context, req *api.GetTokenRequest) (*api.GetTokenResponse, error) {
	iamReq := &iam_api.GetTokenRequest{
		Username: req.Username,
		Password: req.Password,
	}

	resp, err := s.iamClient.GetToken(ctx, iamReq)
	if err != nil {
		return nil, err
	}

	return &api.GetTokenResponse{
		Token: resp.Token,
	}, nil
}

func (s *Service) CreateTicket(ctx context.Context, req *api.CreateTicketRequest) (*api.CreateTicketResponse, error) {
	// check whether the user is valid
	iamReq := &iam_api.GetUserInfoRequest{
		UserId: req.UserId,
	}

	_, err := s.iamClient.GetUserInfo(ctx, iamReq)
	if err != nil {
		return nil, err
	}

	tmReq := &tm_api.CreateTicketRequest{
		UserId:   req.UserId,
		Title:    req.Title,
		Content:  req.Content,
		Priority: req.Priority,
	}

	resp, err := s.tmClient.CreateTicket(ctx, tmReq)
	if err != nil {
		return nil, err
	}

	return &api.CreateTicketResponse{
		Message:  resp.Message,
		TicketId: resp.TicketId,
	}, nil
}

func (s *Service) UpdateTicketStatus(ctx context.Context, req *api.UpdateTicketStatusRequest) (*api.UpdateTicketStatusResponse, error) {
	tmReq := &tm_api.UpdateTicketStatusRequest{
		TicketId: req.TicketId,
		Status:   req.Status,
	}

	// can only update ticket status if the user current status is "pending"
	ticket, err := s.tmClient.GetTicketById(ctx, &tm_api.GetTicketByIdRequest{
		TicketId: req.TicketId,
	})
	if err != nil {
		return nil, err
	}

	if ticket.Status != "pending" {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update ticket status")
	}

	resp, err := s.tmClient.UpdateTicketStatus(ctx, tmReq)
	if err != nil {
		// log error
		return nil, err
	}

	return &api.UpdateTicketStatusResponse{
		Message: resp.Message,
	}, nil
}

func (s *Service) Close(ctx context.Context) {
}

func (s *Service) Ping() error {
	err := s.store.Ping()
	return err
}
