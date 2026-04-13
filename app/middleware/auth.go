package middleware

import (
	"go-production/app/helpers"
	"go-production/app/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	UserContextKey = "user"
)

// Authenticate is a middleware that identifies the user based on the JWT token.
// If no token is provided, it sets an AnonymousUser in the context.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add "Vary: Authorization" header to ensure that responses are cached correctly.
		c.Header("Vary", "Authorization")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Set(UserContextKey, model.AnonymousUser)
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		tokenString := parts[1]
		claims, err := helpers.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Fetch the user from the database
		userModel := model.User{}
		user, err := userModel.GetByID(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		c.Set(UserContextKey, user)
		c.Next()
	}
}

// RequireAuthenticatedUser is a middleware that ensures the user is not anonymous.
func RequireAuthenticatedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get(UserContextKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		user := userInterface.(*model.User)
		if user.IsAnonymous() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you must be logged in to access this resource"})
			return
		}

		c.Next()
	}
}

// GetUser retrieves the user from the context.
func GetUser(c *gin.Context) *model.User {
	user, exists := c.Get(UserContextKey)
	if !exists {
		return model.AnonymousUser
	}
	return user.(*model.User)
}
