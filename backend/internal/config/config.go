package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigtoml"
)

type Config struct {
	Drip      Drip      `toml:"drip" env:"DRIP"`
	Redis     Redis     `toml:"redis" env:"REDIS"`
	Substrate Substrate `toml:"substrate" env:"SUBSTRATE"`

	Discord Discord `toml:"discord" env:"DISCORD"`
	Matrix  Matrix  `toml:"matrix" env:"MATRIX"`
}

type Substrate struct {
	Endpoint     string `toml:"endpoint" env:"ENDPOINT" default:"ws://substrat1e:9944"`
	SeedOrPhrase string `toml:"seed_or_phrase" env:"SEED_OR_PHRASE"`
}

type Redis struct {
	Endpoint string `toml:"endpoint" env:"ENDPOINT"`
}

type Drip struct {
	Cap             float64 `toml:"cap" env:"CAP" default:"0.025"`        // coin per drip
	Delay           int64   `toml:"delay" env:"DELAY" default:"86400000"` // in milliseconds
	NetworkDecimals uint16  `toml:"network_decimals" env:"NETWORK_DECIMALS" default:"12"`
}

type Discord struct {
	Enabled bool   `toml:"enabled" env:"ENABLED" default:"false"`
	Token   string `toml:"token" env:"TOKEN"`
}

type Matrix struct {
	Enabled  bool   `toml:"enabled" env:"ENABLED" default:"false"`
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
