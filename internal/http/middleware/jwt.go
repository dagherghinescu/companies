package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kelseyhightower/envconfig"
)

// JWTConfig holds the configuration needed for the jwt auth implementation.
// Currently it's holding only a secret.
type JWTConfig struct {
	Secret string `envconfig:"SECRET" required:"true"`
}

// EnvConfig loads config from environment variables into HTTPConfig.
func EnvConfig() (*JWTConfig, error) {
	var cfg JWTConfig
	if err := envconfig.Process("JWT", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func JWTMiddleware(cfg *JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		_, err := jwt.Parse(tokenString, func(_ *jwt.Token) (any, error) {
			return []byte(cfg.Secret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Next()
	}
}
