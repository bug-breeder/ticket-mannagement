package service

import (
	"context"

	"go.tekoapis.com/tekone/app/warehouse/tm_service/api"
	health "go.tekoapis.com/tekone/library/grpc/health"

	// "go.tekoapis.com/tekone/app/warehouse/tm_service/config"
	"go.tekoapis.com/tekone/app/warehouse/tm_service/internal/store"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type Service struct {
	log logr.Logger
	// more connector here
	store store.StoreQuerier
	// embedded unimplemented service server
	health.UnimplementedHealthCheckServiceServer
	api.UnimplementedTmServiceServer
}

func NewService(
	logger logr.Logger,
	store store.StoreQuerier,
) *Service {

	return &Service{
		log:   logger,
		store: store,
		// more here

	}
}

func (s *Service) CreateTicket(ctx context.Context, req *api.CreateTicketRequest) (*api.CreateTicketResponse, error) {

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	ticketID, err := s.store.CreateTicket(ctx, store.CreateTicketParams{
		UserID:   userID,
		Title:    req.Title,
		Content:  req.Content,
		Priority: req.Priority,
	})
	if err != nil {
		return nil, err
	}

	return &api.CreateTicketResponse{
		Message:  "Ticket created successfully",
		TicketId: ticketID.String(),
	}, nil
}

func (s *Service) UpdateTicketStatus(ctx context.Context, req *api.UpdateTicketStatusRequest) (*api.UpdateTicketStatusResponse, error) {
	ticketID, err := uuid.Parse(req.TicketId)
	if err != nil {
		return nil, err
	}

	// // status can not be "pending"
	// if req.Status == "pending" {
	// 	return nil, status.Errorf(codes.InvalidArgument, "Status can not be pending")
	// }

	err = s.store.UpdateTicketStatus(ctx, store.UpdateTicketStatusParams{
		TicketID: ticketID,
		Status:   req.Status,
	})

	if err != nil {
		return nil, err
	}

	return &api.UpdateTicketStatusResponse{
		Message: "Ticket status updated successfully",
	}, nil
}

func (s *Service) GetTicketById(ctx context.Context, req *api.GetTicketByIdRequest) (*api.Ticket, error) {
	ticketID, err := uuid.Parse(req.TicketId)
	if err != nil {
		return nil, err
	}

	ticket, err := s.store.GetTicketById(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	return &api.Ticket{
		TicketId: ticket.TicketID.String(),
		UserId:   ticket.UserID.String(),
		Title:    ticket.Title,
		Content:  ticket.Content,
		Priority: ticket.Priority,
		Status:   ticket.Status,
	}, nil
}

func (s *Service) Close(ctx context.Context) {
	s.store.Close()
}

func (s *Service) Ping() error {
	err := s.store.Ping()
	return err
}
