package happening

import "time"

type Config struct {
	Name           string
	ReportURL      string
	SuccessCode    string
	PingURL        string
	FlagHostname   string
	Retries        uint
	RetryDelay     time.Duration
	CollectOutput  bool
	SuppressOutput bool
	Chdir          string
	StoreReport    bool
}

func NewConfig() *Config {
	return &Config{
		Name:        "some event",
		StoreReport: true,
		SuccessCode: "0",
		Retries:     3,
		RetryDelay:  time.Second,
	}
}
