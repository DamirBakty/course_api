package v1

import (
	"github.com/gin-gonic/gin"
	"web/config"
	"web/middleware"
	"web/services"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	app         *config.AppConfig
	service     *services.UserService
	authService *services.AuthService
}

// NewUserHandler claims a new user handler
func NewUserHandler(app *config.AppConfig, service *services.UserService, authService *services.AuthService) *UserHandler {
	return &UserHandler{
		app:         app,
		service:     service,
		authService: authService,
	}
}

// RegisterRoutes registers user api to the router
func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/v1/users")
	userGroup.Use(middleware.AuthMiddleware(h.authService))
	{
		userGroup.POST("/login", h.Claim)
	}
}

// RegisterUser handles POST /api/v1/users/login
// @Summary Register a new user
// @Description Register a new user using Keycloak token
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /users/login [post]
func (h *UserHandler) Claim(c *gin.Context) {
	// Get the claims from the context (set by AuthMiddleware)
	claims, exists := c.Get("claims")
	if !exists {
		middleware.RespondWithError(c, 401, "Authentication required")
		return
	}

	// Convert to KeycloakClaims
	keycloakClaims, ok := claims.(*services.KeycloakClaims)
	if !ok {
		middleware.RespondWithError(c, 500, "Invalid claims type")
		return
	}

	// Register user using token information
	userResponse, err := h.service.ClaimUserUserFromToken(keycloakClaims)
	if err != nil {
		middleware.RespondWithBadRequest(c, err.Error())
		return
	}

	middleware.RespondWithCreated(c, userResponse, "Authorized successfully")
}
