package substrate

import (
	"substrate-faucet/internal/config"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
)

func New(conf config.Substrate) (*gsrpc.SubstrateAPI, error) {
	client, err := gsrpc.NewSubstrateAPI(conf.Endpoint)

	return client, err
}
