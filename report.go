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
		req, err := newHttpClientRequest(http.MethodPut, url, bytes.NewBuffer(jb))
		if err != nil {
			break
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
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
			fmt.Sprintf(
				"giving up connecting %s after %d unsuccessful retries",
				url,
				config.Retries,
			),
		)
	}
	log.Printf("failed, %v.\n", err)
}
