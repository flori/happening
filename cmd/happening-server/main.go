package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	happening "github.com/flori/happening"
)

func basicAuthConfig(config happening.ServerConfig) middleware.BasicAuthConfig {
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

type errorsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func errorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			return c.JSON(http.StatusInternalServerError,
				errorsResponse{Success: false, Message: err.Error()},
			)
		} else {
			return err
		}
	}
}

func main() {
	var config happening.ServerConfig
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(errorHandler)
	if config.HTTP_AUTH != "" {
		fmt.Println("info: Configuring HTTP Auth access control")
		e.Use(middleware.BasicAuthWithConfig(basicAuthConfig(config)))
	}
	api := happening.API{
		DATABASE_NAME: config.DATABASE_NAME,
		POSTGRES_URL:  config.POSTGRES_URL,
		NOTIFIER:      happening.NewNotifier(config),
	}
	api.PrepareDatabase()
	api.SetupCronJobs()
	e.POST("/api/v1/event", api.PostEventHandler)
	e.PUT("/api/v1/event", api.PostEventHandler)
	e.GET("/api/v1/events", api.GetEventsHandler)
	e.POST("/api/v1/check", api.PostCheckHandler)
	e.PUT("/api/v1/check", api.PostCheckHandler)
	e.GET("/api/v1/checks", api.GetChecksHandler)
	e.DELETE("/api/v1/check/:id", api.DeleteCheckHandler)
	e.Logger.Fatal(e.Start(":" + config.PORT))
}
