package discord

import (
	"fmt"
	"substrate-faucet/internal/domain/entity"
	"substrate-faucet/internal/domain/service"

	"github.com/bwmarrin/discordgo"
)

type Service struct {
	session *discordgo.Session
	msgs    chan entity.Message
}

func (s *Service) handleMessage(ss *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == ss.State.User.ID {
		return
	}

	s.msgs <- entity.Message{
		ChannelID: m.ChannelID,
		FromID:    m.Message.Author.ID,
		Body:      []byte(m.Content),
	}
}

// NewDiscordService - constructor for Discord service
func NewDiscordService(token string) (service.UMIService, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("discord new: %+v", err)
	}

	s := &Service{
		session: dg,
		msgs:    make(chan entity.Message),
	}

	dg.AddHandler(s.handleMessage)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	return s, nil
}
