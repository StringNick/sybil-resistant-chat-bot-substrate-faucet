package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"substrate-faucet/internal/config"
	"substrate-faucet/internal/domain/service"
	"substrate-faucet/internal/env/redis"
	"substrate-faucet/internal/env/substrate"
	"substrate-faucet/internal/handler/processor"
	"substrate-faucet/internal/service/drip"
	"substrate-faucet/internal/service/umi/discord"
	"substrate-faucet/internal/service/umi/matrix"
	"syscall"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"go.uber.org/zap"
)

func prepareLogger() {
	logger, err := zap.NewDevelopment()

	//logger, err := zapdriver.NewDevelopment()
	if err != nil {
		log.Fatal()
	}
	zap.ReplaceGlobals(logger)

}

// registerServices - register all services
func registerServices(cfg *config.Config) []service.UMIService {
	var services []service.UMIService

	if cfg.Matrix.Enabled {
		matrix, err := matrix.NewMatrix(context.Background(), matrix.NewMatrixParams{
			DeviceID: cfg.Matrix.DeviceID,
			Host:     cfg.Matrix.Host,
			Username: cfg.Matrix.Username,
			Password: cfg.Matrix.Password,
		})
		if err != nil {
			panic(err)
		}

		services = append(services, matrix)
	}

	if cfg.Discord.Enabled {
		ds, err := discord.NewDiscordService(cfg.Discord.Token)
		if err != nil {
			panic(err)
		}

		services = append(services, ds)
	}

	if len(services) == 0 {
		panic("no services provided")
	}

	return services
}

func main() {
	// prepare zap logger
	prepareLogger()

	// reading config
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	// creating redis client
	rdb, err := redis.NewRedis(cfg.Redis)
	if err != nil {
		panic(err)
	}

	// creating substrate client
	sc, err := substrate.New(cfg.Substrate)
	if err != nil {
		panic(err)
	}

	// transferFrom account
	transferFrom, err := signature.KeyringPairFromSecret(cfg.Substrate.SeedOrPhrase, 0)
	if err != nil {
		panic(err)
	}

	// creating drip service
	dripSvc, err := drip.New(drip.Params{
		Rdb:                 rdb,
		SubstrateClient:     sc,
		SubstrateTransferer: transferFrom,

		Cap:      cfg.Drip.Cap,
		CapDelay: cfg.Drip.Delay,
	})
	if err != nil {
		panic(err)
	}

	// register all umi services
	services := registerServices(cfg)

	zap.L().Debug("registered services", zap.Any("services_count", len(services)))

	procHandler := processor.NewHandler(processor.HandlerParams{
		Config: cfg,

		UMIServices: services,

		DripService: dripSvc,
	})
	if err != nil {
		panic(err)
	}

	zap.L().Debug("starting services processor...")

	// starting processor
	if err = procHandler.Start(); err != nil {
		panic(err)
	}

	zap.L().Debug("successfully started processor")
	zap.L().Debug("handling command and waiting for interrupt signal...")

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// waiting interrupt signal
	<-sigs

	zap.L().Debug("got interrupt signal, stopping processor...")

	procHandler.Stop()

	zap.L().Debug("successfully stopped processor")
}
