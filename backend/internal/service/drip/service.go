package drip

import (
	"context"
	"fmt"
	"substrate-faucet/internal/domain/service"
	"substrate-faucet/internal/env/substrate"
	"time"

	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

var fmtLastDripKey = "last_drip_%s"

// GetLastDrip - get last drip time, if not exist ErrLastDripNotFound returns
func (s *Service) GetLastDrip(address string) (time.Time, error) {
	val, err := s.rdb.Get(context.Background(), fmt.Sprintf(fmtLastDripKey, address)).Int64()
	if err != nil {
		if err == redis.Nil {
			return time.Time{}, service.ErrLastDripNotFound
		}

		return time.Time{}, err
	}

	return time.Unix(val, 0), nil
}

// UpdateLastDrip - updating last drip time
func (s *Service) UpdateLastDrip(address string) error {
	_, err := s.GetLastDrip(address)
	if err == nil {
		return service.ErrDripAlreadyExist
	} else if err != service.ErrLastDripNotFound {
		return err
	}

	tm := time.Now().Add(time.Duration(s.capDelay) * time.Second)

	err = s.rdb.Set(context.Background(), fmt.Sprintf(fmtLastDripKey, address), tm.Unix(), time.Second*time.Duration(s.capDelay)).Err()
	if err != nil {
		return err
	}

	// trying to send tx in substrate
	hash, err := substrate.MakeATransfer(s.substrateClient, s.substrateTransferer, address, 1)
	if err != nil {
		// deleting key, because we have an error
		defer s.rdb.Del(context.Background(), fmt.Sprintf(fmtLastDripKey, address))
		return err
	}

	zap.L().Debug("substrate tx sent", zap.String("address", address), zap.String("hash", hash))
	return nil
}
