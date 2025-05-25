package rest

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// JwtCustomClaims defines the JWT payload structure.
type JwtCustomClaims struct {
	UserID any    `json:"user_id"`
	App    string `json:"app"`
	jwt.RegisteredClaims
}

// GenerateJwtToken creates a JWT token for any type of userID (int, string, etc).
func GenerateJwtToken(userID any, app string) (string, error) {
	claims := JwtCustomClaims{
		UserID: userID,
		App:    app,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8766 * time.Hour)), // ~1 year
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	return signed, err
}

// Restricted returns JWT middleware for Echo framework.
func Restricted() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("JWT_SECRET")),
		TokenLookup: "query:token,header:Authorization",
	})
}

// ParseUserIDFromToken extracts the user_id from JWT claims.
func ParseUserIDFromToken(c echo.Context) (any, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["user_id"], nil
}
