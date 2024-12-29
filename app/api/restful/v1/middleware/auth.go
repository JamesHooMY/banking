package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"banking/domain"
	"banking/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware is a middleware to protect routes with JWT authentication
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Authorization token required"})
			c.Abort()
			return
		}

		// Bearer token format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Authorization token format is Bearer {token}"})
			c.Abort()
			return
		}

		// Parse the token
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("authedUserId", claims.UserID)
		c.Set("isAdmin", claims.IsAdmin)
		c.Set("email", claims.Email)
		c.Next()
	}
}

func APIKeyAuthMiddleware(
	authService domain.IAuthService,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve API key and Secret key from headers
		key := c.GetHeader("X-API-Key")
		secretKey := c.GetHeader("X-Secret-Key")

		userIDStr := c.GetHeader("X-User-Id")
		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil || userID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid User ID"})
			c.Abort()
			return
		}

		if key == "" || secretKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "API Key and Secret Key are required"})
			c.Abort()
			return
		}

		// Use gin.Context for propagation
		ctx := c.Request.Context()

		// check if the API key and secret key are valid
		if err := authService.APIKeyConfirmation(ctx, uint(userID), key, secretKey); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Invalid API Key or Secret Key"})
			c.Abort()
			return
		}

		c.Set("authedUserId", uint(userID))
		c.Set("apiKey", key)
		c.Set("secretKey", secretKey)

		// Continue processing the request
		c.Next()
	}
}
