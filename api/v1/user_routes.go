package v1

import (
	"github.com/gin-gonic/gin"
	"web/config"
	"web/middleware"
	"web/schemas"
	"web/services"
)

type UserHandler struct {
	app         *config.AppConfig
	service     *services.UserService
	authService *services.AuthService
}

func NewUserHandler(app *config.AppConfig, service *services.UserService, authService *services.AuthService) *UserHandler {
	return &UserHandler{
		app:         app,
		service:     service,
		authService: authService,
	}
}

// RegisterRoutes registers user api to the router
func (h *UserHandler) RegisterRoutes(router *gin.Engine) {
	// Public routes (no authentication required)
	publicGroup := router.Group("/api/v1/auth")
	{
		publicGroup.POST("/login", h.Login)
		publicGroup.POST("/refresh", h.RefreshToken)
	}

	// Protected routes (authentication required)
	protectedGroup := router.Group("/api/v1/users")
	protectedGroup.Use(middleware.AuthMiddleware(h.authService))
	{
		// Admin-only routes
		adminGroup := protectedGroup.Group("/admin")
		adminGroup.Use(middleware.RequireRole(h.authService, "admin"))
		{
			adminGroup.POST("/create", h.AdminCreateUser)
		}
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

// Login handles POST /api/v1/auth/login
// @Summary Login with username and password
// @Description Authenticate with Keycloak and get a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body schemas.LoginRequest true "Login credentials"
// @Success 200 {object} schemas.LoginResponse "Login successful"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Authentication failed"
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var loginRequest schemas.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		middleware.RespondWithBadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	// Authenticate with Keycloak
	loginResponse, err := h.authService.Login(loginRequest.Username, loginRequest.Password, *h.service)
	if err != nil {
		middleware.RespondWithError(c, 401, "Authentication failed: "+err.Error())
		return
	}
	middleware.RespondWithSuccess(c, loginResponse, "Login successful")
}

// RefreshToken handles POST /api/v1/auth/refresh
// @Summary Refresh an access token
// @Description Refresh an access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param refresh body schemas.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} schemas.LoginResponse "Token refreshed successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Token refresh failed"
// @Router /auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var refreshRequest schemas.RefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		middleware.RespondWithBadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	// Refresh the token
	loginResponse, err := h.authService.RefreshToken(refreshRequest.RefreshToken)
	if err != nil {
		middleware.RespondWithError(c, 401, "Token refresh failed: "+err.Error())
		return
	}

	middleware.RespondWithSuccess(c, loginResponse, "Token refreshed successfully")
}

// AdminCreateUser handles POST /api/v1/users/admin/create
// @Summary Create a new user (Admin only)
// @Description Create a new user in Keycloak and in the local database (Admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body schemas.AdminCreateUserRequest true "User data"
// @Success 201 {object} schemas.UserResponse "User created successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Router /users/admin/create [post]
func (h *UserHandler) AdminCreateUser(c *gin.Context) {
	var userRequest schemas.AdminCreateUserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		middleware.RespondWithBadRequest(c, "Invalid request body: "+err.Error())
		return
	}

	// Create user
	userResponse, err := h.service.AdminCreateUser(userRequest, h.authService)
	if err != nil {
		middleware.RespondWithBadRequest(c, err.Error())
		return
	}

	middleware.RespondWithCreated(c, userResponse, "User created successfully")
}
