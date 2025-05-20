package v1

import (
	"github.com/gin-gonic/gin"
	"web/config"
	"web/middleware"
	"web/schemas"
	"web/services"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	app         *config.AppConfig
	service     *services.UserService
	authService *services.AuthService
}

// NewUserHandler creates a new user handler
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
	{
		// Only users with ROLE_ADMIN can register new users
		userGroup.POST("/register", middleware.RequireRole(h.authService, "ROLE_ADMIN"), h.RegisterUser)
	}
}

// RegisterUser handles POST /api/v1/users/register
// @Summary Register a new user
// @Description Register a new user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param user body schemas.RegisterUserRequest true "User data"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - requires admin role"
// @Router /users/register [post]
// @example request - example payload
//
//	{
//	  "username": "johndoe",
//	  "email": "john.doe@example.com",
//	  "password": "password123",
//	  "roles": "ROLE_USER"
//	}
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var userRequest schemas.RegisterUserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		middleware.RespondWithBadRequest(c, "Invalid request body")
		return
	}

	userResponse, err := h.service.RegisterUser(userRequest)
	if err != nil {
		middleware.RespondWithBadRequest(c, err.Error())
		return
	}

	middleware.RespondWithCreated(c, userResponse, "User registered successfully")
}
