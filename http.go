package happening

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"net/http"
)

func newHttpClient() *http.Client {
	rootCAs := x509.NewCertPool()
	if certPool, _ := x509.SystemCertPool(); certPool != nil {
		rootCAs = certPool
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAs,
			},
		},
	}

	return client
}

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
