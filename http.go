package happening

import (
	"bytes"
	"net/http"
)

func newHttpClientRequest(method, url string, buffer *bytes.Buffer) (*http.Request, error) {
	if buffer == nil {
		buffer = &bytes.Buffer{}
	}
	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "happening")
	return req, err
}
