package substrate

import (
	"substrate-faucet/internal/domain/service"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
)

type Service struct {
	api *gsrpc.SubstrateAPI
}

type Params struct {
	API *gsrpc.SubstrateAPI
}

func New(params Params) (service.SubstrateService, error) {
	return &Service{
		api: params.API,
	}, nil
}
