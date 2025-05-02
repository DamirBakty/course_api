package v1

import (
	"net/http"
	"strconv"
	"web/config"
	"web/models"
	"web/schemas"
	"web/services"

	"github.com/gin-gonic/gin"
)

// ChapterHandler handles HTTP requests for chapters
type ChapterHandler struct {
	app     *config.AppConfig
	service *services.ChapterService
}

// NewChapterHandler creates a new chapter handler
func NewChapterHandler(app *config.AppConfig, service *services.ChapterService) *ChapterHandler {
	return &ChapterHandler{
		app:     app,
		service: service,
	}
}

// RegisterRoutes registers chapter api to the router
func (h *ChapterHandler) RegisterRoutes(router *gin.Engine) {
	courseGroup := router.Group("/api/v1/courses")
	{
		chapterGroup := courseGroup.Group("/:id/chapters")
		{
			chapterGroup.GET("", h.GetAllChapters)
			chapterGroup.GET("/:chapterId", h.GetChapterByID)
			chapterGroup.POST("", h.CreateChapter)
			chapterGroup.PUT("/:chapterId", h.UpdateChapter)
			chapterGroup.DELETE("/:chapterId", h.DeleteChapter)
		}
	}
}

// GetAllChapters handles GET /api/v1/courses/:id/chapters
// @Summary Get all chapters for a course
// @Description Get a list of all chapters for a specific course
// @Tags chapters
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} map[string]interface{} "Returns a list of chapters"
// @Failure 400 {object} map[string]interface{} "Invalid course ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/chapters [get]
func (h *ChapterHandler) GetAllChapters(c *gin.Context) {
	courseIdStr := c.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid course ID",
		})
		return
	}

	chapters, err := h.service.GetChaptersByCourseID(uint(courseId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  chapters,
	})
}

// GetChapterByID handles GET /api/v1/courses/:id/chapters/:chapterId
// @Summary Get a chapter by ID
// @Description Get a chapter by its ID
// @Tags chapters
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param chapterId path int true "Chapter ID"
// @Success 200 {object} map[string]interface{} "Returns the chapter"
// @Failure 400 {object} map[string]interface{} "Invalid chapter ID"
// @Failure 404 {object} map[string]interface{} "Chapter not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/chapters/{chapterId} [get]
func (h *ChapterHandler) GetChapterByID(c *gin.Context) {
	courseIdStr := c.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid course ID",
		})
		return
	}

	idStr := c.Param("chapterId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid chapter ID",
		})
		return
	}

	chapter, err := h.service.GetChapterByID(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "chapter not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Verify that the chapter belongs to the specified course
	if chapter.CourseID != uint(courseId) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Chapter not found for this course",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  chapter,
	})
}

// CreateChapter handles POST /api/v1/courses/:id/chapters
// @Summary Create a new chapter
// @Description Create a new chapter with the provided data
// @Tags chapters
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param chapter body schemas.ChapterRequest true "Chapter data"
// @Success 201 {object} map[string]interface{} "Chapter created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Router /courses/{id}/chapters [post]
// @example request - example payload
//
//	{
//	  "name": "Chapter 1: Getting Started",
//	  "description": "Introduction to the course material",
//	  "order": 1
//	}
func (h *ChapterHandler) CreateChapter(c *gin.Context) {
	courseIdStr := c.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid course ID",
		})
		return
	}
	var chapterRequest schemas.ChapterRequest
	if err := c.ShouldBindJSON(&chapterRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}

	// Set the courseId from the URL path
	chapter := models.Chapter{
		CourseID:    uint(courseId),
		Name:        chapterRequest.Name,
		Description: chapterRequest.Description,
		Order:       chapterRequest.Order,
	}

	id, err := h.service.CreateChapter(chapter)
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
		"message": "Chapter created successfully",
	})
}

// UpdateChapter handles PUT /api/v1/courses/:id/chapters/:chapterId
// @Summary Update a chapter
// @Description Update a chapter with the provided data
// @Tags chapters
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param chapterId path int true "Chapter ID"
// @Param chapter body models.Chapter true "Chapter data"
// @Success 200 {object} map[string]interface{} "Chapter updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 404 {object} map[string]interface{} "Chapter not found"
// @Router /courses/{id}/chapters/{chapterId} [put]
func (h *ChapterHandler) UpdateChapter(c *gin.Context) {
	courseIdStr := c.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid course ID",
		})
		return
	}

	idStr := c.Param("chapterId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid chapter ID",
		})
		return
	}

	var chapterRequest schemas.ChapterRequest
	if err := c.ShouldBindJSON(&chapterRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request body",
		})
		return
	}

	chapter := models.Chapter{
		Name:        chapterRequest.Name,
		Description: chapterRequest.Description,
		Order:       chapterRequest.Order,
	}
	chapter.ID = uint(id)
	chapter.CourseID = uint(courseId)
	err = h.service.UpdateChapter(chapter)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "chapter not found or no changes made" {
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
		"message": "Chapter updated successfully",
	})
}

// DeleteChapter handles DELETE /api/v1/courses/:id/chapters/:chapterId
// @Summary Delete a chapter
// @Description Delete a chapter by its ID
// @Tags chapters
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param chapterId path int true "Chapter ID"
// @Success 200 {object} map[string]interface{} "Chapter deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid chapter ID"
// @Failure 404 {object} map[string]interface{} "Chapter not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/chapters/{chapterId} [delete]
func (h *ChapterHandler) DeleteChapter(c *gin.Context) {
	courseIdStr := c.Param("id")
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid course ID",
		})
		return
	}

	idStr := c.Param("chapterId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid chapter ID",
		})
		return
	}

	// Verify that the chapter belongs to the specified course
	chapter, err := h.service.GetChapterByID(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "chapter not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	if chapter.CourseID != uint(courseId) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Chapter not found for this course",
		})
		return
	}

	err = h.service.DeleteChapter(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "chapter not found" {
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
		"message": "Chapter deleted successfully",
	})
}
