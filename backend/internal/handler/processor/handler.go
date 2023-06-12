package processor

import (
	"context"
	"fmt"
	"substrate-faucet/internal/domain/entity"
	"substrate-faucet/internal/domain/service"
	"time"

	"go.uber.org/zap"
)

// processMsg - processing message from umi service
func (h *Handler) processMsg(s service.UMIService, msg *entity.Message) {
	zap.L().Debug("msg from some umi service", zap.Any("msg", msg))
}

// processUmiService - processing umi service for messages
func (h *Handler) processUmiService(s service.UMIService) {
	defer h.wg.Done()

	for {
		ctx, cancel := context.WithTimeout(h.ctx, time.Second)

		msg, err := s.RcvMsg(ctx)
		cancel()

		switch err {
		case context.Canceled:
			// so we out of processing, context was canceled
			zap.L().Debug("timeout exceed")
			return
		case context.DeadlineExceeded:
			// so we timeouted, let's try again
			continue
		default:
			if err != nil {
				// so we got some custom error
				zap.L().Debug("undefined error", zap.Error(err))
				continue
			}

			h.processMsg(s, msg)
		}
	}
}

// Start - non blocking start handler
func (h *Handler) Start() error {
	// starting all provided umi services
	for _, s := range h.umiServices {
		if err := s.Start(); err != nil {
			return fmt.Errorf("umi service start: %+v", err)
		}
	}

	h.wg.Add(len(h.umiServices))

	// start handling messages from all umi services in background
	for _, s := range h.umiServices {
		go h.processUmiService(s)
	}

	return nil
}

func (h *Handler) Stop() {
	h.cancel()

	// waiting when all background services stop handling messages
	h.wg.Wait()

	// closing connections with all umi services
	for _, s := range h.umiServices {
		if err := s.Stop(); err != nil {
			zap.L().Debug("umi service stop", zap.Error(err))
		}
	}
}
