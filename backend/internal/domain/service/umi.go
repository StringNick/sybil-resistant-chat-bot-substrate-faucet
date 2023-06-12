package service

import (
	"context"
	"substrate-faucet/internal/domain/entity"
)

// UMIService - its some provider of chat platform, like Slack, Telegram, Discord, etc.
type UMIService interface {
	// Start - starting handling messages (non-blocking)
	Start() error
	// Stop - stopping handling messages (blocking)
	Stop() error

	// RcvMsg - blocking read from chat platform
	RcvMsg(ctx context.Context) (*entity.Message, error)
	// SndMsg - blocking write to chat platform
	SndMsg(ctx context.Context, channelID, replyTo string, body []byte) error
}
