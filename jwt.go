package happening

import (
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JwtAuth(config ServerConfig) echojwt.Config {
	return echojwt.Config{
		Skipper: func(c echo.Context) bool {
			path := c.Path()
			method := c.Request().Method
			return path == "/api/v1/event" && (method == "POST" || method == "PUT")
		},
		SigningKey: []byte(config.SIGNING_SECRET),
	}
}

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type AuthPair struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func JwtLoginWithConfig(config ServerConfig) func(echo.Context) error {
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
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
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
