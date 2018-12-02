package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	h "github.com/flori/happening"
)

type Config struct {
	PORT          string `default:"8080"`
	DATABASE_NAME string `default:"happening"`
	POSTGRES_URL  string `default:"postgresql://flori@localhost:5432/%s?sslmode=disable"`
	HTTP_REALM    string `default:"happening"`
	HTTP_AUTH     string
}

func basicAuthConfig(config Config) middleware.BasicAuthConfig {
	return middleware.BasicAuthConfig{
		Realm: config.HTTP_REALM,
		Skipper: func(c echo.Context) bool {
			return c.Request().Method != "GET"
		},
		Validator: func(username, password string, c echo.Context) (bool, error) {
			httpAuth := strings.Split(config.HTTP_AUTH, ":")
			if username == httpAuth[0] && password == httpAuth[1] {
				return true, nil
			}
			return false, nil
		},
	}
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	if config.HTTP_AUTH != "" {
		fmt.Println("info: Configuring HTTP Auth access control")
		e.Use(middleware.BasicAuthWithConfig(basicAuthConfig(config)))
	}
	api := h.API{DATABASE_NAME: config.DATABASE_NAME, POSTGRES_URL: config.POSTGRES_URL}
	api.PrepareDatabase()
	e.POST("/api/v1/event", api.PostEventHandler)
	e.GET("/api/v1/events", api.GetEventsHandler)
	e.Logger.Fatal(e.Start(":" + config.PORT))
}
