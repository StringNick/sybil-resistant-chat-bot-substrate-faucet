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
	if len(msg.Body) > len("/request ") && string(msg.Body[:len("/request ")]) == "/request " {
		addr := string(msg.Body[len("/request "):])

		zap.L().Debug("request for drip", zap.String("addr", addr))

		drip, err := h.dripService.GetLastDrip(addr)
		if err != nil && err != service.ErrLastDripNotFound {
			// error when get last drip (unexpected)
			s.SndMsg(context.Background(), msg.ChannelID, msg.FromID, []byte("Something happend wrong, try again later!"))
			zap.L().Error("something wrong when get last drip", zap.String("address", addr), zap.Error(err))
			return
		} else if err == service.ErrLastDripNotFound {
			// so drip not found we can take again
			if err = h.dripService.UpdateLastDrip(addr); err != nil {
				if err == service.ErrWrongAddress {
					// display wrong address for drip
					s.SndMsg(context.Background(), msg.ChannelID, msg.FromID, []byte("Sorry ur address is not SS58 encoded!"))
				} else {
					s.SndMsg(context.Background(), msg.ChannelID, msg.FromID, []byte("Something happend wrong, try again later!"))
					zap.L().Error("something wrong when update last drip", zap.String("address", addr), zap.Error(err))
				}
				return
			}

			zap.L().Debug("successfully updated last drip", zap.String("address", addr))
			err = s.SndMsg(context.Background(), msg.ChannelID, msg.FromID, []byte("Successfully sent drip!"))
			if err != nil {
				zap.L().Error("something wrong when send message", zap.String("address", addr), zap.Error(err))
				return
			}
			return
		}

		// drip already exist we can print to chat that we can't take again
		err = s.SndMsg(context.Background(), msg.ChannelID, msg.FromID,
			[]byte(fmt.Sprintf("You can't take drip again, last drip was at %s. U can do it after: %s", drip.Format(time.RFC3339), drip.Add(time.Duration(h.cfg.Drip.Delay)*time.Millisecond).Format(time.RFC3339))))
		if err != nil {
			zap.L().Error("something wrong when send message", zap.String("address", addr), zap.Error(err))
			return
		}
	} else {
		zap.L().Debug("unexpected msg in chat")
	}
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
