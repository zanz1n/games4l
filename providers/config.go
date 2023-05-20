package providers

import (
	"os"

	"github.com/games4l/telemetry-service/logger"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
)

type Config struct {
	Port         uint16 `json:"port,omitempty"`
	WebhookSig   string `json:"webhook_sig,omitempty"`
	MongoUri     string `json:"mongo_uri,omitempty"`
	MongoDbName  string `json:"mongo_db_name,omitempty"`
	RoutePrefix  string `json:"route_prefix"`
	ProjectEpoch int64  `json:"project_epoch"`
}

var (
	config   Config
	validate = validator.New()
)

func parseJSON(buf []byte) {
	err := json.Unmarshal(buf, &config)

	if err != nil {
		logger.Fatal(err)
	}

	err = validate.Struct(config)

	if err != nil {
		logger.Fatal("invalid config")
	}
}

func GetConfig() *Config {
	return &config
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
