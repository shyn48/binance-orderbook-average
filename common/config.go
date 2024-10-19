package common

import (
	"sync"

	load_env "github.com/amovah/load-env"
	"github.com/joho/godotenv"
)

type BinanaceConnectionConfig struct {
	Endpoint    string `env:"name=BINANCE_ENDPOINT"`
	Symbol      string `env:"name=BINANCE_SYMBOL"`
	Depth       int    `env:"name=BINANCE_DEPTH"`
	UpdateSpeed int    `env:"name=BINANCE_UPDATE_SPEED"`
}

type Config struct {
	Port     int `env:"name=PORT,default=8080"`
	LogLevel int `env:"name=LOG_LEVEL,default=5"`
	Binance  BinanaceConnectionConfig
}

var (
	config     Config
	configErr  error
	configOnce sync.Once
)

func GetConfig() (Config, error) {
	configOnce.Do(func() {
		godotenv.Load()
		configErr = load_env.LoadEnv(&config)
	})

	return config, configErr
}
