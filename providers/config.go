package providers

import (
	"os"

	"github.com/games4l/telemetria/logger"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
)

type Config struct {
	Port       uint8  `json:"port"`
	WebhookSig string `json:"webhook_sig"`
}

var (
	config   *Config
	validate = validator.New()
)

func parseJSON(buf []byte) {
	json.Unmarshal(buf, &config)

	err := validate.Struct(config)

	if err != nil {
		logger.Fatal("invalid config")
	}
}

func GetConfig() *Config {
	return config
}

func AcquireFromFile(path string) {
	cfgBuf, err := os.ReadFile(path)

	if err != nil {
		logger.Fatal("config could not be found")
	}

	parseJSON(cfgBuf)
}

func AcquireFromEnv() {
	if envJson := os.Getenv("APP_CONFIG"); envJson != "" {
		parseJSON([]byte(envJson))
	} else {
		logger.Fatal("config env not found")
	}
}
