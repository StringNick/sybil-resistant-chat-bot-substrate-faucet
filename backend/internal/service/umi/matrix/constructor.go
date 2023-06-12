package matrix

import (
	"context"
	"fmt"
	"substrate-faucet/internal/domain/entity"
	"substrate-faucet/internal/domain/service"
	"sync"
	"time"

	"go.uber.org/zap"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type Service struct {
	msgs chan entity.Message
	cli  *mautrix.Client

	wg      *sync.WaitGroup
	stopMu  sync.Mutex
	stopped bool

	mach *crypto.OlmMachine
}

func (s *Service) handleMessage(source mautrix.EventSource, evt *event.Event) {
	s.msgs <- entity.Message{
		ChannelID: string(evt.RoomID),
		FromID:    string(evt.Sender),
	}
}

func (s *Service) getUserIDs(roomID id.RoomID) ([]id.UserID, error) {
	members, err := s.cli.JoinedMembers(roomID)
	if err != nil {
		return nil, err
	}
	userIDs := make([]id.UserID, len(members.Joined))
	i := 0
	for userID := range members.Joined {
		userIDs[i] = userID
		i++
	}
	return userIDs, nil
}

func (s *Service) sendEncrypted(ctx context.Context, roomID id.RoomID, text string) error {
	content := event.MessageEventContent{
		MsgType: "m.text",
		Body:    text,
		Format:  event.FormatHTML,
	}

	encrypted, err := s.mach.EncryptMegolmEvent(context.Background(), roomID, event.EventMessage, content)
	// These three errors mean we have to make a new Megolm session
	if crypto.IsShareError(err) {
		userIds, err := s.getUserIDs(roomID)
		if err != nil {
			return fmt.Errorf("getUserIds error: %+v", err)
		}

		err = s.mach.ShareGroupSession(context.Background(), roomID, userIds)
		if err != nil {
			return err
		}
		encrypted, err = s.mach.EncryptMegolmEvent(ctx, roomID, event.EventMessage, content)
		if err != nil {
			return fmt.Errorf("s.mach.EncryptMegolmEvent: %+v", err)
		}

	} else if err != nil {
		return fmt.Errorf("encrypt megolm event err: %+v", err)
	}

	_, err = s.cli.SendMessageEvent(roomID, event.EventEncrypted, encrypted)
	if err != nil {
		return err
	}

	return nil
}

type NewMatrixParams struct {
	DeviceID string
	Host     string
	Username string
	Password string
}

// NewMatrix - creating matrix service instance of UMI
func NewMatrix(ctx context.Context, p NewMatrixParams) (service.UMIService, error) {
	start := time.Now().UnixNano() / 1_000_000
	cli, err := mautrix.NewClient(p.Host, "", "")
	if err != nil {
		return nil, err
	}

	// Log in to get access token and device ID.
	_, err = cli.Login(&mautrix.ReqLogin{
		Type: mautrix.AuthTypePassword,
		Identifier: mautrix.UserIdentifier{
			Type: mautrix.IdentifierTypeUser,
			User: p.Username,
		},
		Password:                 p.Password,
		InitialDeviceDisplayName: p.DeviceID,
		//DeviceID:                 id.DeviceID(deviceID),
		StoreCredentials: true,
	})
	if err != nil {
		return nil, err
	}

	// Create a store for the e2ee keys. In real apps, use NewSQLCryptoStore instead of NewGobStore.
	cryptoStore := crypto.NewMemoryStore(nil)

	//log := zerolog.New(nil)

	mach := crypto.NewOlmMachine(cli, nil, cryptoStore, &fakeStateStore{})
	// Load data from the crypto store
	err = mach.Load()
	if err != nil {
		return nil, err
	}

	svc := &Service{
		msgs: make(chan entity.Message),
		cli:  cli,
		mach: mach,

		wg: &sync.WaitGroup{},
	}
	// Hook up the OlmMachine into the Matrix client so it receives e2ee keys and other such things.
	syncer := cli.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnSync(func(resp *mautrix.RespSync, since string) bool {
		mach.ProcessSyncResponse(resp, since)
		return true
	})
	syncer.OnEventType(event.StateMember, func(source mautrix.EventSource, evt *event.Event) {
		mach.HandleMemberEvent(source, evt)
	})
	// Listen to encrypted messages
	syncer.OnEventType(event.EventEncrypted, func(source mautrix.EventSource, evt *event.Event) {
		if evt.Timestamp < start {
			// Ignore events from before the program started
			return
		}
		decrypted, err := mach.DecryptMegolmEvent(context.Background(), evt)
		if err != nil {
			zap.L().Debug("Failed to decrypt", zap.Error(err))
		} else {
			message, isMessage := decrypted.Content.Parsed.(*event.MessageEventContent)

			// checking if it message and sender is not bot
			if isMessage && decrypted.Sender != cli.UserID {
				svc.msgs <- entity.Message{
					ChannelID: decrypted.RoomID.String(),
					FromID:    decrypted.Sender.String(),
					Body:      []byte(message.Body),
				}
			}
		}
	})

	return svc, nil
}
