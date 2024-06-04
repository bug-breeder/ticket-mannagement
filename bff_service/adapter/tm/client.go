package tm

import (
	"context"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"

	"go.tekoapis.com/tekone/app/warehouse/tm_service/api"
)

type ClientAdapter interface {
	CreateTicket(ctx context.Context, req *api.CreateTicketRequest) (*api.CreateTicketResponse, error)
	UpdateTicketStatus(ctx context.Context, req *api.UpdateTicketStatusRequest) (*api.UpdateTicketStatusResponse, error)
	GetTicketById(ctx context.Context, req *api.GetTicketByIdRequest) (*api.Ticket, error)
}

type client struct {
	client api.TmServiceClient
	log    logr.Logger
}

func New(log logr.Logger, addr string) (ClientAdapter, error) {
	tmConn, err := grpc.DialContext(context.Background(), addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	tmClient := api.NewTmServiceClient(tmConn)
	return &client{
		client: tmClient,
		log:    log,
	}, nil
}

func (c *client) CreateTicket(ctx context.Context, req *api.CreateTicketRequest) (*api.CreateTicketResponse, error) {
	return c.client.CreateTicket(ctx, req)
}

func (c *client) UpdateTicketStatus(ctx context.Context, req *api.UpdateTicketStatusRequest) (*api.UpdateTicketStatusResponse, error) {
	return c.client.UpdateTicketStatus(ctx, req)
}

func (c *client) GetTicketById(ctx context.Context, req *api.GetTicketByIdRequest) (*api.Ticket, error) {
	return c.client.GetTicketById(ctx, req)
}
