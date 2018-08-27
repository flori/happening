package happening

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Ping(config *Config) {
	url := config.PingURL
	var err error
	for i := uint(0); i < config.Retries; i++ {
		log.Printf("Pinging %s…", url)
		resp, err := http.Get(url)
		if err != nil {
			time.Sleep(config.RetryDelay)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode < 400 {
			log.Println("succeeded.")
			return
		}
		time.Sleep(config.RetryDelay)
	}
	if err == nil {
		err = errors.New(
			fmt.Sprintf("giving up pinging %s after %d unsuccessful retries", url, config.Retries))
	}
	log.Printf("failed, %v.\n", err)
}
