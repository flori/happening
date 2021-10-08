package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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
				log.Printf("Access for %s was granted.", username)
				return true, nil
			}
			log.Printf("Access for %s was denied!", username)
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

	api := happening.API{
		DATABASE_NAME: config.DATABASE_NAME,
		POSTGRES_URL:  config.POSTGRES_URL,
		NOTIFIER:      happening.NewNotifier(config),
	}
	api.PrepareDatabase()
	api.SetupCronJobs()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORS())
	e.Use(errorHandler)

	g := e.Group("/api/v1")
	if config.HTTP_AUTH == "" {
		log.Fatal("Need HTTP_AUTH configuration to start server.")
	} else {
		fmt.Println("info:dmi Configuring HTTP Auth access control")
		g.Use(middleware.BasicAuthWithConfig(basicAuthConfig(config)))
	}
	if config.NOTIFIER_KIND != "" {
		log.Printf("Using notifier for %s", config.NOTIFIER_KIND)
	} else {
		log.Printf("Notifier disabled.")
	}

	// Inserting events
	g.POST("/event", api.PostEventHandler)
	g.PUT("/event", api.PostEventHandler)

	// Events
	g.GET("/events", api.GetEventsHandler)
	g.GET("/event/:id", api.GetEventHandler)

	// Checks
	g.POST("/check", api.PostCheckHandler)
	g.PUT("/check", api.PostCheckHandler)
	g.PATCH("/check/:id", api.PatchCheckHandler)
	g.GET("/checks", api.GetChecksHandler)
	g.DELETE("/check/:id", api.DeleteCheckHandler)
	g.GET("/check/:id", api.GetCheckHandler)
	g.GET("/check/by-name/:name", api.GetCheckByNameHandler)

	// HTML
	e.GET("/*", happening.PackrHandler(config))

	e.Logger.Fatal(e.Start(":" + config.PORT))
}
