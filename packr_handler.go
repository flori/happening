package happening

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr"
	"github.com/labstack/echo"
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

func PackrHandler(config ServerConfig) echo.HandlerFunc {
	createEnv(config)
	box := packr.NewBox(filepath.Join(config.WEBUI_DIR, "build"))
	fileServer := http.FileServer(box)
	wrapHandler := func(h http.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			if path != "" && box.Has(path) {
				h.ServeHTTP(c.Response(), c.Request())
				return nil
			} else {
				if path == "" {
					path = "/search"
				}
				return c.Redirect(http.StatusTemporaryRedirect, "/?p="+path)
			}
		}
	}
	return wrapHandler(fileServer)
}
