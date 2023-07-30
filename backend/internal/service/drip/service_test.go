package drip

import (
	"context"
	"reflect"
	"substrate-faucet/internal/domain/service"
	"testing"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/redis/go-redis/v9"
	"github.com/vedhavyas/go-subkey"
)

type mockSubstrateService struct{}

func (m *mockSubstrateService) Transfer(sig signature.KeyringPair, addr string, amnt uint64) (string, error) {
	_, _, err := subkey.SS58Decode(addr)
	if err != nil {
		return "", service.ErrWrongAddress
	}

	return "0x1234567890", nil
}

func TestWrongAddress(t *testing.T) {
	api := &mockSubstrateService{}

	_, err := api.Transfer(signature.TestKeyringPairAlice, "7FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty", 1)
	if err != service.ErrWrongAddress {
		t.Fatal(err)
	}
}

func TestMakeATransfer(t *testing.T) {
	api := &mockSubstrateService{}

	keyring, err := signature.KeyringPairFromSecret("//Alice", 42)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(keyring, signature.TestKeyringPairAlice) {
		t.Logf("got: %+v, except: %+v\n", keyring, signature.TestKeyringPairAlice)
		t.Fatal("keyring is not equal to signature.TestKeyringPairAlice")
	}

	hash, err := api.Transfer(keyring, "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty", 1)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("substrate tx sent, hash: %s\n", hash)
}

type redisTest struct {
	d map[string]interface{}
}

func newRedisTest() *redisTest {
	return &redisTest{
		d: make(map[string]interface{}),
	}
}

func (r *redisTest) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	r.d[key] = value
	return redis.NewStatusResult("OK", nil)
}

func (r *redisTest) Get(ctx context.Context, key string) *redis.StringCmd {
	v, ok := r.d[key]
	if !ok {
		return redis.NewStringResult("", redis.Nil)
	}

	return redis.NewStringResult(v.(string), nil)
}

func (r *redisTest) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	for _, key := range keys {
		delete(r.d, key)
	}
	return redis.NewIntResult(1, nil)
}

func TestService_Test(t *testing.T) {
	rt := newRedisTest()

	api := &mockSubstrateService{}

	svc, err := New(Params{
		Rdb:              rt,
		SubstrateService: api,
		SubstrateTransferer: signature.KeyringPair{
			URI:       "",
			Address:   "",
			PublicKey: []byte{},
		},
		Cap:      1.0,
		CapDelay: 10000,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.GetLastDrip("test")
	if err != service.ErrLastDripNotFound {
		t.Fatal(err)
	}

	err = svc.UpdateLastDrip("test")
	if err != service.ErrWrongAddress {
		t.Fatal(err)
	}

	err = svc.UpdateLastDrip("5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty")
	if err != nil {
		t.Fatal(err)
	}
}
