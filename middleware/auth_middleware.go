package middleware

import (
	"net/http"
	"strings"

	"gadget-marketplace/pkg/utils"

	"github.com/labstack/echo/v4"
)

func JWTAuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "Authorization header is required",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "Authorization header format must be 'Bearer <token>'",
				})
			}

			tokenString := parts[1]
			claims, err := utils.ValidateJWTToken(tokenString, jwtSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"error": "Invalid or expired token: " + err.Error(),
				})
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}
