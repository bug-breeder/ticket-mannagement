package pubsub

import (
	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
	"go.tekoapis.com/tekone/app/warehouse/bff_service/config"
)

type PubSub interface {
	pubsubInterface
}

type pubsub struct {
	log logr.Logger
	cfg *config.Config
	*pubsubCore
}

func NewPubSub(log logr.Logger, cfg *config.Config) PubSub {
	ps := &pubsub{
		log: log,
		cfg: cfg,
	}
	ps.pubsubCore = newPubsubCore(cfg)
	return ps
}

func (ps *pubsub) Start(cliCtx *cli.Context) error {
	errChan := make(chan error)
	processes := []interface{}{}
	go ps.startSubscriber(errChan, processes...)
	return <-errChan
}
