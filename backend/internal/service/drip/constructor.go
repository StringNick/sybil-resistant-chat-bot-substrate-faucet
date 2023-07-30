package drip

import (
	"context"
	"substrate-faucet/internal/domain/service"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type Params struct {
	Rdb                 RedisClient
	SubstrateTransferer signature.KeyringPair

	SubstrateService service.SubstrateService

	Cap             float64
	CapDelay        int64
	NetworkDecimals uint16
}

type Service struct {
	rdb                 RedisClient
	substrateService    service.SubstrateService
	substrateTransferer signature.KeyringPair

	cap             float64
	capDelay        int64
	networkDecimals uint16
}

func New(params Params) (service.DripService, error) {
	return &Service{
		rdb:      params.Rdb,
		cap:      params.Cap,
		capDelay: params.CapDelay,

		substrateService:    params.SubstrateService,
		substrateTransferer: params.SubstrateTransferer,
		networkDecimals:     params.NetworkDecimals,
	}, nil
}
