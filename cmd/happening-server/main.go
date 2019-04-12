package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	happening "github.com/flori/happening"
)

func basicAuthConfig(config happening.ServerConfig) middleware.BasicAuthConfig {
	return middleware.BasicAuthConfig{
		Realm: config.HTTP_REALM,
		Skipper: func(c echo.Context) bool {
			m := c.Request().Method
			return c.Path() == "/api/v1/event" && (m == "POST" || m == "PUT")
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
			if httpErr, ok := err.(*echo.HTTPError); ok {
				return c.JSON(
					httpErr.Code,
					errorsResponse{Success: false, Message: err.Error()},
				)
			} else {
				return c.JSON(
					http.StatusInternalServerError,
					errorsResponse{Success: false, Message: err.Error()},
				)
			}
		} else {
			return err
		}
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	var config happening.ServerConfig
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	//e.Use(middleware.Logger())
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
	log.Printf("Using notifier for %s", config.NOTIFIER_KIND)
	api.PrepareDatabase()
	api.SetupCronJobs()
	// Events
	e.POST("/api/v1/event", api.PostEventHandler)
	e.PUT("/api/v1/event", api.PostEventHandler)
	e.GET("/api/v1/events", api.GetEventsHandler)
	// Checks
	e.POST("/api/v1/check", api.PostCheckHandler)
	e.PUT("/api/v1/check", api.PostCheckHandler)
	e.PATCH("/api/v1/check/:id", api.PatchCheckHandler)
	e.GET("/api/v1/checks", api.GetChecksHandler)
	e.DELETE("/api/v1/check/:id", api.DeleteCheckHandler)
	e.GET("/api/v1/check/:id", api.GetCheckHandler)
	e.GET("/api/v1/check/by-name/:name", api.GetCheckByNameHandler)
	//
	e.Logger.Fatal(e.Start(":" + config.PORT))
}
