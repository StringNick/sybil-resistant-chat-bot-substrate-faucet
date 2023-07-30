package service

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
)

var ErrWrongAddress = fmt.Errorf("wrong address")

type SubstrateService interface {
	Transfer(signature.KeyringPair, string, uint64) (string, error)
}
