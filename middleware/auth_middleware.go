package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redukasquad/be-reduka/database/migrations"
	"github.com/redukasquad/be-reduka/modules/users"
	"github.com/redukasquad/be-reduka/packages/utils"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed("Unauthorized", "Missing authorization header", nil))
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed("Unauthorized", "Invalid token", nil))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed("Unauthorized", "Invalid token claims", nil))
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed("Unauthorized", "Invalid user ID in token", nil))
			return
		}

		c.Set("user_id", int(userIDFloat))
		c.Next()
	}
}

func RequireAuthorization(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed("Unauthorized", "User ID not found in context", nil))
			return
		}

		userRepo := users.NewRepository(migrations.GetDB())
		user, err := userRepo.FindByID(userID.(int))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed("Unauthorized", "User not found", nil))
			return
		}

		hasRole := false
		for _, role := range allowedRoles {
			if strings.EqualFold(user.Role, role) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.BuildResponseFailed("Forbidden", fmt.Sprintf("User role '%s' is not authorized", user.Role), nil))
			return
		}

		c.Next()
	}
}
