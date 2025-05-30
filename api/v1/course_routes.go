package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"web/config"
	"web/middleware"
	"web/schemas"
	"web/services"
)

// CourseHandler handles HTTP requests for courses
type CourseHandler struct {
	app            *config.AppConfig
	service        *services.CourseService
	chapterService *services.ChapterService
	authService    *services.AuthService
}

// NewCourseHandler creates a new course handler
func NewCourseHandler(app *config.AppConfig, service *services.CourseService, chapterService *services.ChapterService, authService *services.AuthService) *CourseHandler {
	return &CourseHandler{
		app:            app,
		service:        service,
		chapterService: chapterService,
		authService:    authService,
	}
}

// RegisterRoutes registers course api to the router
func (h *CourseHandler) RegisterRoutes(router *gin.Engine) {
	courseGroup := router.Group("/api/v1/courses")
	courseGroup.Use(middleware.AuthMiddleware(h.authService))
	{
		courseGroup.GET("", h.GetAllCourses)
		courseGroup.GET("/:id", h.GetCourseByID)
		courseGroup.POST("", h.CreateCourse)
		courseGroup.PUT("/:id", h.UpdateCourse)
		courseGroup.DELETE("/:id", h.DeleteCourse)
	}
}

// GetAllCourses handles GET /api/courses
// @Summary Get all courses
// @Description Get a list of all courses
// @Tags courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Returns a list of courses"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses [get]
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	courseResponses, err := h.service.GetAllCourses()
	if err != nil {
		middleware.RespondWithInternalServerError(c, err.Error())
		return
	}

	middleware.RespondWithSuccess(c, courseResponses, "")
}

// GetCourseByID handles GET /api/courses/:id
// @Summary Get a course by ID
// @Description Get a course by its ID
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Success 200 {object} map[string]interface{} "Returns the course"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id} [get]
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		middleware.RespondWithBadRequest(c, "Invalid course ID")
		return
	}

	courseResponse, err := h.service.GetCourseByIDWithChapterCount(uint(id))
	if err != nil {
		if err.Error() == "course not found" {
			middleware.RespondWithNotFound(c, err.Error())
		} else {
			middleware.RespondWithInternalServerError(c, err.Error())
		}
		return
	}

	middleware.RespondWithSuccess(c, courseResponse, "")
}

// CreateCourse handles POST /api/courses
// @Summary Create a new course
// @Description Create a new course with the provided data
// @Tags courses
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param course body schemas.CreateCourseRequest true "Course data"
// @Success 201 {object} map[string]interface{} "Course created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Router /courses [post]
// @example request - example payload
//
//	{
//	  "name": "Introduction to Go Programming",
//	  "description": "Learn the basics of Go programming language"
//	}
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var courseRequest schemas.CreateCourseRequest
	if err := c.ShouldBindJSON(&courseRequest); err != nil {
		middleware.RespondWithBadRequest(c, "Invalid request body")
		return
	}

	sub, exists := c.Get("sub")
	if !exists {
		middleware.RespondWithBadRequest(c, "User not found in context")
		return
	}

	// Get the user by sub
	user, err := h.authService.GetUserBySub(sub.(string))
	if err != nil {
		middleware.RespondWithBadRequest(c, "Failed to get user: "+err.Error())
		return
	}

	// Set the created_by field
	userID := user.ID
	courseRequest.CreatedBy = &userID

	courseResponse, err := h.service.CreateCourse(courseRequest)
	if err != nil {
		middleware.RespondWithBadRequest(c, err.Error())
		return
	}

	middleware.RespondWithCreated(c, courseResponse, "Course created successfully")
}

// UpdateCourse handles PUT /api/courses/:id
// @Summary Update a course
// @Description Update a course with the provided data
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param course body schemas.UpdateCourseRequest true "Course data"
// @Success 200 {object} map[string]interface{} "Course updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Router /courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		middleware.RespondWithBadRequest(c, "Invalid course ID")
		return
	}

	var courseRequest schemas.UpdateCourseRequest

	if err := c.ShouldBindJSON(&courseRequest); err != nil {
		middleware.RespondWithBadRequest(c, "Invalid request body")
		return
	}

	course, err := h.service.GetCourseByID(uint(id))
	if err != nil {
		if err.Error() == "course not found" {
			middleware.RespondWithNotFound(c, err.Error())
		} else {
			middleware.RespondWithInternalServerError(c, err.Error())
		}
		return
	}

	courseResponse, err := h.service.UpdateCourse(course, courseRequest)
	if err != nil {
		if err.Error() == "course not found or no changes made" {
			middleware.RespondWithNotFound(c, err.Error())
		} else {
			middleware.RespondWithBadRequest(c, err.Error())
		}
		return
	}

	middleware.RespondWithSuccess(c, courseResponse, "Course updated successfully")
}

// DeleteCourse handles DELETE /api/courses/:id
// @Summary Delete a course
// @Description Delete a course by its ID
// @Tags courses
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Success 200 {object} map[string]interface{} "Course deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 404 {object} map[string]interface{} "Course not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id} [delete]
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		middleware.RespondWithBadRequest(c, "Invalid course ID")
		return
	}

	err = h.service.DeleteCourse(uint(id))
	if err != nil {
		if err.Error() == "course not found" {
			middleware.RespondWithNotFound(c, err.Error())
		} else {
			middleware.RespondWithInternalServerError(c, err.Error())
		}
		return
	}

	middleware.RespondWithSuccess(c, nil, "Course deleted successfully")
}
