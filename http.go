package happening

import (
	"bytes"
	"log"
	"net/http"
)

func newHttpClientRequest(method, url string, buffer *bytes.Buffer) *http.Request {
	if buffer == nil {
		buffer = &bytes.Buffer{}
	}
	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		log.Printf("failed, %v.\n", err)
		return nil
	}
	req.Header.Set("User-Agent", "happening")
	return req
}
