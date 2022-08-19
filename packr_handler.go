package happening

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr/v2"
	"github.com/labstack/echo/v4"
)

type env struct {
	HAPPENING_SERVER_URL string `json:"HAPPENING_SERVER_URL"`
}

func createEnv(config ServerConfig) {
	vars := env{
		HAPPENING_SERVER_URL: config.HAPPENING_SERVER_URL,
	}

	bytes, err := json.MarshalIndent(&vars, "", "  ")
	if err != nil {
		log.Panic(err)
	}

	outputPath := filepath.Join(config.WEBUI_DIR, "build", "Env.js")
	file, err := os.Create(outputPath)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Env = %s", bytes)
	_, err = fmt.Fprintf(file, "window.Env = %s;", bytes)
	if err != nil {
		log.Panic(err)
	}
}

func redirectURL(path string, oldURL *url.URL) string {
	query := oldURL.Query()
	path = path + "?" + query.Encode()
	oldURL.Path = "/"
	newQuery := url.Values{}
	newQuery.Set("p", path)
	oldURL.RawQuery = newQuery.Encode()
	return oldURL.String()
}

func PackrHandler(config ServerConfig) echo.HandlerFunc {
	createEnv(config)
	box := packr.New("myBox", filepath.Join(config.WEBUI_DIR, "build"))
	fileServer := http.FileServer(box)
	wrapHandler := func(h http.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			oldURL := c.Request().URL
			path := oldURL.Path
			if path != "" && box.Has(path) {
				log.Printf("Serving %s from box.", path)
				h.ServeHTTP(c.Response(), c.Request())
				return nil
			} else {
				if path == "" {
					path = "/search"
				}
				log.Printf("Redirecting URL from %s to %s.", oldURL, path)
				return c.Redirect(http.StatusTemporaryRedirect, redirectURL(path, oldURL))
			}
		}
	}
	return wrapHandler(fileServer)
}
