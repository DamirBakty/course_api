package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web/services"
)

// AuthMiddleware creates a middleware that validates JWT tokens from Keycloak
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		token, err := authService.ExtractToken(c.Request)
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// Validate the token locally
		claims, err := authService.ValidateToken(token)
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
			c.Abort()
			return
		}

		// Check with Keycloak if the token is still valid (not revoked or blocked)
		active, err := authService.IntrospectToken(token)
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, "Token validation failed: "+err.Error())
			c.Abort()
			return
		}
		if !active {
			RespondWithError(c, http.StatusUnauthorized, "Token is no longer active")
			c.Abort()
			return
		}

		// Validate the session (check if user exists in the database)
		if c.Request.URL.Path != "/api/v1/users/login" { // Skip session validation for login endpoint
			valid, err := authService.ValidateSession(claims.Sub)
			if err != nil {
				RespondWithError(c, http.StatusInternalServerError, "Session validation error: "+err.Error())
				c.Abort()
				return
			}
			if !valid {
				userService := services.NewUserService(authService.GetUserRepo())

				_, createErr := userService.ClaimUserUserFromToken(claims)
				if createErr != nil {
					RespondWithError(c, http.StatusInternalServerError, "Failed to create user: "+createErr.Error())
					c.Abort()
					return
				}

				_, getUserErr := authService.GetUserBySub(claims.Sub)
				if getUserErr != nil {
					RespondWithError(c, http.StatusInternalServerError, "Failed to get created user: "+getUserErr.Error())
					c.Abort()
					return
				}
			}
		}
		user, err := authService.GetUserBySub(claims.Sub)
		if err != nil {
			RespondWithError(c, http.StatusUnauthorized, "Invalid session: user not found")
			c.Abort()
			return
		}
		// Store the claims in the context for later use
		c.Set("claims", claims)
		c.Set("user_id", claims.Subject)
		c.Set("user", user)
		c.Set("username", claims.PreferredUsername)
		c.Set("email", claims.Email)
		c.Set("sub", claims.Sub)

		c.Next()
	}
}

// RequireRole creates a middleware that requires a specific role
func RequireRole(authService *services.AuthService, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the claims from the context
		claims, exists := c.Get("claims")
		if !exists {
			RespondWithError(c, http.StatusUnauthorized, "Authentication required")
			c.Abort()
			return
		}

		// Check if the user has the required role
		keycloakClaims, ok := claims.(*services.KeycloakClaims)
		if !ok {
			RespondWithError(c, http.StatusInternalServerError, "Invalid claims type")
			c.Abort()
			return
		}

		if !authService.HasRole(keycloakClaims, role) {
			RespondWithError(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}
