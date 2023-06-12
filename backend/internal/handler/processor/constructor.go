package processor

import (
	"context"
	"substrate-faucet/internal/config"
	"substrate-faucet/internal/domain/service"
	"sync"
)

type Handler struct {
	cfg         *config.Config
	umiServices []service.UMIService
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

type HandlerParams struct {
	// TODO: redis, matrix and discord services
	Config *config.Config

	UMIServices []service.UMIService
}

func NewHandler(params HandlerParams) *Handler {
	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	return &Handler{
		cfg:         params.Config,
		umiServices: params.UMIServices,
		wg:          sync.WaitGroup{},
		ctx:         ctx,
		cancel:      cancel,
	}
}
