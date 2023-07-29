package drip

import (
	"substrate-faucet/internal/domain/service"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"

	"github.com/redis/go-redis/v9"
)

type Params struct {
	Rdb                 *redis.Client
	SubstrateClient     *gsrpc.SubstrateAPI
	SubstrateTransferer signature.KeyringPair

	Cap      float64
	CapDelay int64
}

type Service struct {
	rdb                 *redis.Client
	substrateClient     *gsrpc.SubstrateAPI
	substrateTransferer signature.KeyringPair

	cap      float64
	capDelay int64
}

func New(params Params) (service.DripService, error) {
	return &Service{
		rdb:      params.Rdb,
		cap:      params.Cap,
		capDelay: params.CapDelay,

		substrateClient:     params.SubstrateClient,
		substrateTransferer: params.SubstrateTransferer,
	}, nil
}
