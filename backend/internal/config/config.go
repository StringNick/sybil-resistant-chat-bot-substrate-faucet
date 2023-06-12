package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigtoml"
)

type Config struct {
	Drip      Drip      `toml:"drip" env:"DRIP"`
	Redis     Redis     `toml:"redis" env:"REDIS"`
	Substrate Substrate `toml:"substrate" env:"SUBSTRATE"`
	Faucet    Faucet    `toml:"faucet" env:"FAUCET"`

	Discord Discord `toml:"discord" env:"DISCORD"`
	Matrix  Matrix  `toml:"matrix" env:"MATRIX"`
}

type Faucet struct {
	Secret string `toml:"secret" env:"SECRET"`
}

type Substrate struct {
	Endpoint string `toml:"endpoint" env:"ENDPOINT"`
}

type Redis struct {
	Endpoint string `toml:"endpoint" env:"ENDPOINT"`
}

type Drip struct {
	Cap   float64 `toml:"cap" env:"CAP"`     // coin per drip
	Delay int64   `toml:"delay" env:"DELAY"` // in milliseconds
}

type Discord struct {
	Enabled bool   `toml:"enabled" env:"ENABLED"`
	Token   string `toml:"token" env:"TOKEN"`
}

type Matrix struct {
	Enabled  bool   `toml:"enabled" env:"ENABLED"`
	DeviceID string `toml:"device_id" env:"DEVICE_ID"`
	Host     string `toml:"host" env:"HOST"`
	Username string `toml:"username" env:"USERNAME"`
	Password string `toml:"password" env:"PASSWORD"`
}

func New() (*Config, error) {
	var cfg Config

	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		AllowUnknownEnvs: true,
		SkipFlags:        true,
		EnvPrefix:        "",
		Files:            []string{"./config/config.toml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".toml": aconfigtoml.New(),
		},
	})

	if err := loader.Load(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
