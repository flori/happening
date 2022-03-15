package happening

import (
	"fmt"
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

func switchDatabase(postgresURL string, databaseName string) string {
	pu, err := url.Parse(postgresURL)
	if err != nil {
		log.Panic(err)
	}
	pu.Path = fmt.Sprintf("/%s", databaseName)
	return pu.String()
}
