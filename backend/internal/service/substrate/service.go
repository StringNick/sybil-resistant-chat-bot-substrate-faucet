package substrate

import (
	"fmt"
	"substrate-faucet/internal/domain/service"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey"
)

func (s *Service) Transfer(transferFrom signature.KeyringPair, transferToSS58 string, amount uint64) (string, error) {
	// parse ss58 address
	_, toAccountId, err := subkey.SS58Decode(transferToSS58)
	if err != nil {
		return "", service.ErrWrongAddress
	}

	meta, err := s.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return "", fmt.Errorf("get meta data: %+v", err)
	}

	// Create a call
	transferTo, err := types.NewMultiAddressFromAccountID(toAccountId)
	if err != nil {
		return "", fmt.Errorf("new accountId: %+v", err)
	}

	c, err := types.NewCall(meta, "Balances.transfer", transferTo, types.NewUCompactFromUInt(amount))
	if err != nil {
		return "", fmt.Errorf("new call: %+v", err)
	}

	// Create the extrinsic
	ext := types.NewExtrinsic(c)

	genesisHash, err := s.api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return "", fmt.Errorf("get block hash: %+v", err)
	}

	rv, err := s.api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return "", fmt.Errorf("get runtime version: %+v", err)
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", transferFrom.PublicKey)
	if err != nil {
		return "", fmt.Errorf("storage key: %+v", err)
	}

	var accountInfo types.AccountInfo
	ok, err := s.api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		return "", fmt.Errorf("get storage latest: %+v", err)
	}

	nonce := uint32(accountInfo.Nonce)

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction using keyring pair
	err = ext.Sign(transferFrom, o)
	if err != nil {
		return "", fmt.Errorf("sign: %+v", err)
	}

	// Send the extrinsic
	hash, err := s.api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		return "", fmt.Errorf("submit extrinsic: %+v", err)
	}

	return hash.Hex(), nil
}
