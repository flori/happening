package happening

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Name           string `default:"SomeEvent"`
	Context        string `default:"default" envconfig:"CONTEXT"`
	ReportURL      string `envconfig:"HAPPENING_REPORT_URL"`
	SuccessCode    string `default:"0"`
	PingURL        string
	Hostname       string
	Retries        uint          `default:"3"`
	RetryDelay     time.Duration `default:"1s"`
	CollectOutput  bool
	SuppressOutput bool
	Chdir          string
	StoreReport    bool `default:"true"`
	Started        string
	Duration       time.Duration
	Output         string
}

func NewConfig() *Config {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
