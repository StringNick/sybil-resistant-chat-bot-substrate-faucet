package discord

import (
	"context"
	"fmt"
	"substrate-faucet/internal/domain/entity"
)

// Start - starting handling messages (non-blocking)
func (s Service) Start() error {
	err := s.session.Open()
	if err != nil {
		return fmt.Errorf("session open: %+v", err)
	}

	return nil
}

// Stop - stopping handling messages (blocking)
func (s Service) Stop() error {
	return s.session.Close()
}

// RcvMsg - blocking read from chat platform
func (s Service) RcvMsg(ctx context.Context) (*entity.Message, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case v := <-s.msgs:
		return &v, nil
	}
}

// SndMsg - blocking write to chat platform
func (s Service) SndMsg(ctx context.Context, channelID, replyTo string, data []byte) error {
	_, err := s.session.ChannelMessageSend(channelID, fmt.Sprintf("<@!%s> %s", replyTo, data))
	return err
}
