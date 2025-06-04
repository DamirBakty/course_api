package v1

import (
	"net/http"
	"strconv"
	"web/config"
	"web/middleware"
	"web/schemas"
	"web/services"

	"github.com/gin-gonic/gin"
)

// LessonHandler handles HTTP requests for lessons
type LessonHandler struct {
	app         *config.AppConfig
	service     *services.LessonService
	authService *services.AuthService
}

// NewLessonHandler creates a new lesson handler
func NewLessonHandler(app *config.AppConfig, service *services.LessonService, authService *services.AuthService) *LessonHandler {
	return &LessonHandler{
		app:         app,
		service:     service,
		authService: authService,
	}
}

// RegisterRoutes registers lesson api to the router
func (h *LessonHandler) RegisterRoutes(router *gin.Engine) {
	courseGroup := router.Group("/api/v1/courses")
	{
		chapterGroup := courseGroup.Group("/:id/chapters")
		{
			lessonGroup := chapterGroup.Group("/:chapterId/lessons")
			lessonGroup.Use(middleware.AuthMiddleware(h.authService))
			{
				lessonGroup.GET("", h.GetAllLessons)
				lessonGroup.GET("/:lessonId", h.GetLessonByID)
				lessonGroup.POST("", h.CreateLesson)
				lessonGroup.PUT("/:lessonId", h.UpdateLesson)
				lessonGroup.DELETE("/:lessonId", h.DeleteLesson)
			}
		}
	}
}

// GetAllLessons handles GET /api/v1/courses/:id/chapters/:chapterId/lessons
// @Summary Get all lessons for a chapter
// @Description Get a list of all lessons for a specific chapter
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param chapterId path int true "Chapter ID"
// @Success 200 {object} map[string]interface{} "Returns a list of lessons"
// @Failure 400 {object} map[string]interface{} "Invalid chapter ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/chapters/{chapterId}/lessons [get]
func (h *LessonHandler) GetAllLessons(c *gin.Context) {
	chapterIdStr := c.Param("chapterId")
	courseIdStr := c.Param("id")
	courseID, err := strconv.ParseUint(courseIdStr, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid course ID",
		})
		return
	}

	chapterId, err := strconv.ParseUint(chapterIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid chapter ID",
		})
		return
	}

	lessons, err := h.service.GetLessonsByChapterID(uint(courseID), uint(chapterId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  lessons,
	})
}

// GetLessonByID handles GET /api/v1/courses/:id/chapters/:chapterId/lessons/:lessonId
// @Summary Get a lesson by ID
// @Description Get a lesson by its ID
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param chapterId path int true "Chapter ID"
// @Param lessonId path int true "Lesson ID"
// @Success 200 {object} map[string]interface{} "Returns the lesson"
// @Failure 400 {object} map[string]interface{} "Invalid lesson ID"
// @Failure 404 {object} map[string]interface{} "Lesson not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/chapters/{chapterId}/lessons/{lessonId} [get]
func (h *LessonHandler) GetLessonByID(c *gin.Context) {
	courseIdStr := c.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid course ID",
		})
		return
	}
	chapterIdStr := c.Param("chapterId")
	chapterId, err := strconv.ParseUint(chapterIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid chapter ID",
		})
		return
	}

	idStr := c.Param("lessonId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid lesson ID",
		})
		return
	}

	lesson, err := h.service.GetLessonByID(uint(courseId), uint(chapterId), uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "lesson not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  lesson,
	})
}

// CreateLesson handles POST /api/v1/courses/:id/chapters/:chapterId/lessons
// @Summary Create a new lesson
// @Description Create a new lesson with the provided data
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param chapterId path int true "Chapter ID"
// @Param lesson body schemas.LessonRequest true "Lesson data"
// @Success 201 {object} map[string]interface{} "Lesson created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Router /courses/{id}/chapters/{chapterId}/lessons [post]
// @example request - example payload
//
//	{
//	  "name": "Lesson 1: Introduction",
//	  "description": "Overview of the chapter content",
//	  "content": "This lesson covers the basic concepts of the chapter.",
//	  "order": 1
//	}
func (h *LessonHandler) CreateLesson(c *gin.Context) {
	chapterIdStr := c.Param("chapterId")
	chapterId, err := strconv.ParseUint(chapterIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid chapter ID",
		})
		return
	}
	courseIdStr := c.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid course ID",
		})
		return
	}

	var lessonRequest schemas.LessonRequest
	if err := c.ShouldBindJSON(&lessonRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}

	sub, exists := c.Get("sub")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "User not found in context",
		})
		return
	}

	// Get the user by sub
	user, err := h.authService.GetUserBySub(sub.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Failed to get user: " + err.Error(),
		})
		return
	}

	// Set the created_by field
	userID := user.ID
	lessonRequest.CreatedBy = &userID

	id, err := h.service.CreateLesson(lessonRequest, uint(courseId), uint(chapterId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error": false,
		"data": gin.H{
			"id": id,
		},
		"message": "Lesson created successfully",
	})
}

// UpdateLesson handles PUT /api/v1/courses/:id/chapters/:chapterId/lessons/:lessonId
// @Summary Update a lesson
// @Description Update a lesson with the provided data
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param chapterId path int true "Chapter ID"
// @Param lessonId path int true "Lesson ID"
// @Param lesson body schemas.LessonRequest true "Lesson data"
// @Success 200 {object} map[string]interface{} "Lesson updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 404 {object} map[string]interface{} "Lesson not found"
// @Router /courses/{id}/chapters/{chapterId}/lessons/{lessonId} [put]
func (h *LessonHandler) UpdateLesson(c *gin.Context) {
	courseIdStr := c.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid course ID",
		})
		return
	}
	chapterIdStr := c.Param("chapterId")
	chapterId, err := strconv.ParseUint(chapterIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid chapter ID",
		})
		return
	}

	idStr := c.Param("lessonId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid lesson ID",
		})
		return
	}

	var lessonRequest schemas.LessonRequest
	if err := c.ShouldBindJSON(&lessonRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}

	err = h.service.UpdateLesson(uint(courseId), uint(chapterId), uint(id), lessonRequest)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "lesson not found or no changes made" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Lesson updated successfully",
	})
}

// DeleteLesson handles DELETE /api/v1/courses/:id/chapters/:chapterId/lessons/:lessonId
// @Summary Delete a lesson
// @Description Delete a lesson by its ID
// @Tags lessons
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param chapterId path int true "Chapter ID"
// @Param lessonId path int true "Lesson ID"
// @Success 200 {object} map[string]interface{} "Lesson deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid lesson ID"
// @Failure 404 {object} map[string]interface{} "Lesson not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/chapters/{chapterId}/lessons/{lessonId} [delete]
func (h *LessonHandler) DeleteLesson(c *gin.Context) {
	courseIdStr := c.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid course ID",
		})
		return
	}
	chapterIdStr := c.Param("chapterId")
	chapterId, err := strconv.ParseUint(chapterIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid chapter ID",
		})
		return
	}

	idStr := c.Param("lessonId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid lesson ID",
		})
		return
	}

	err = h.service.DeleteLesson(uint(courseId), uint(chapterId), uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "lesson not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Lesson deleted successfully",
	})
}
