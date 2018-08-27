package happening

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func EventToJSON(event *Event) []byte {
	jsonBuffer, err := json.Marshal(event)
	if err != nil {
		log.Fatal(err)
	}
	return jsonBuffer
}

func SendEvent(event *Event, config *Config) {
	url := config.ReportURL
	var err error
	jb := EventToJSON(event)
	for i := uint(0); i < config.Retries; i++ {
		log.Printf("Sending event \"%s\" to %s…", event.Name, url)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jb))
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
