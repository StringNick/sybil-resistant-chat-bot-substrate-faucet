package substrate

import (
	"reflect"
	"substrate-faucet/internal/config"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
)

func TestMakeATransferWrongAddress(t *testing.T) {
	_, err := MakeATransfer(nil, signature.KeyringPair{}, "sssss", 1)
	if err != ErrWrongAddress {
		t.Fatal("err is not ErrWrongAddress", err)
	}
}

func TestMakeATransfer(t *testing.T) {
	api, err := New(config.Substrate{
		Endpoint: "ws://localhost:9944",
	})
	if err != nil {
		t.Fatal(err)
	}

	keyring, err := signature.KeyringPairFromSecret("//Alice", 42)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(keyring, signature.TestKeyringPairAlice) {
		t.Logf("got: %+v, except: %+v\n", keyring, signature.TestKeyringPairAlice)
		t.Fatal("keyring is not equal to signature.TestKeyringPairAlice")
	}

	hash, err := MakeATransfer(api, keyring, "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty", 1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("substrate tx sent, hash: %s\n", hash)
}
