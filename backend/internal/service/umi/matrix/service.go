package matrix

import (
	"context"
	"fmt"
	"substrate-faucet/internal/domain/entity"

	"go.uber.org/zap"
	"maunium.net/go/mautrix/id"
)

// Start - starting handling messages (non-blocking)
func (s *Service) Start() error {
	s.wg.Add(1)

	// Start long polling in the background
	go func() {
		defer s.wg.Done()

		err := s.cli.SyncWithContext(context.Background())
		if err != nil {
			zap.L().Debug("sync error", zap.Error(err))
			panic(err)
		}
	}()
	return nil
}

// Stop - stopping service handling
func (s *Service) Stop() error {
	s.stopMu.Lock()
	if s.stopped {
		s.stopMu.Unlock()
		return nil
	}

	s.stopped = true
	s.stopMu.Unlock()

	s.cli.StopSync()
	s.wg.Wait()
	return nil
}

func (s *Service) RcvMsg(ctx context.Context) (*entity.Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case v := <-s.msgs:
		return &v, nil
	}
}

func (s *Service) SndMsg(ctx context.Context, channelID, replyTo string, data []byte) error {
	return s.sendEncrypted(ctx, id.RoomID(channelID), fmt.Sprintf("%s, %s", replyTo, data))
}
