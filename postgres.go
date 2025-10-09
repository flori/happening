package happening

import (
	"log"
	"net/url"
	"path"
)

func deriveDatabaseName(postgresURL string) string {
	pu, err := url.Parse(postgresURL)
	if err != nil {
		log.Panic(err)
	}
	return path.Base(pu.Path)
}
