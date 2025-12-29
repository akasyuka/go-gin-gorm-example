package security

import (
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// InitJWKS инициализирует JWKS ключи Keycloak
func InitJWKS(jwksURL string) (*keyfunc.JWKS, error) {
	options := keyfunc.Options{
		RefreshInterval: time.Hour,
	}
	return keyfunc.Get(jwksURL, options)
}

// JWTMiddleware проверяет JWT токен из Authorization header
func JWTMiddleware(jwks *keyfunc.JWKS) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		tokenStr := parts[1]

		// Разбор токена с использованием JWKS
		token, err := jwt.Parse(tokenStr, jwks.Keyfunc)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Сохраняем claims в контекст Gin
		c.Set("claims", token.Claims)
		c.Next()
	}
}
