package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	happening "github.com/flori/happening"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func jwtAuth(config happening.ServerConfig) middleware.JWTConfig {
	return middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			path := c.Path()
			method := c.Request().Method
			return path == "/api/v1/event" && (method == "POST" || method == "PUT")
		},
		SigningKey: []byte(config.SIGNING_SECRET),
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

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type AuthPair struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func jwtLoginWithConfig(config happening.ServerConfig) func(echo.Context) error {
	return func(c echo.Context) error {
		authPair := new(AuthPair)
		if err := c.Bind(authPair); err != nil {
			return err
		}

		httpAuth := strings.Split(config.HTTP_AUTH, ":")

		// Throws unauthorized error
		if authPair.Username != httpAuth[0] || authPair.Password != httpAuth[1] {
			return echo.ErrUnauthorized
		}

		// Set custom claims
		claims := &jwtCustomClaims{
			"Admin",
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(), // 1 week
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(config.SIGNING_SECRET))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
	}
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
	if config.HTTP_AUTH == "" || config.SIGNING_SECRET == "" {
		log.Fatal("Need HTTP_AUTH and SIGNING_SECRET configuration to start server.")
	} else {
		fmt.Println("info:dmi Configuring JWT authentication")
		g.Use(middleware.JWTWithConfig(jwtAuth(config)))
	}
	e.POST("/jwt", jwtLoginWithConfig(config))
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
