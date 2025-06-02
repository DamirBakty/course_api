package v1

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"web/config"
	"web/middleware"
	"web/services"

	"github.com/gin-gonic/gin"
)

// AttachmentHandler handles HTTP requests for attachments
type AttachmentHandler struct {
	app         *config.AppConfig
	service     *services.AttachmentService
	authService *services.AuthService
}

// NewAttachmentHandler creates a new attachment handler
func NewAttachmentHandler(app *config.AppConfig, service *services.AttachmentService, authService *services.AuthService) *AttachmentHandler {
	return &AttachmentHandler{
		app:         app,
		service:     service,
		authService: authService,
	}
}

// RegisterRoutes registers attachment api to the router
func (h *AttachmentHandler) RegisterRoutes(router *gin.Engine) {
	// Routes for attachments following the hierarchical structure
	courseGroup := router.Group("/api/v1/courses")
	{
		chapterGroup := courseGroup.Group("/:id/chapters")
		{
			lessonGroup := chapterGroup.Group("/:chapterId/lessons")
			{
				attachmentGroup := lessonGroup.Group("/:lessonId/attachments")
				attachmentGroup.Use(middleware.AuthMiddleware(h.authService))
				{
					// GET all attachments for a lesson
					attachmentGroup.GET("", h.GetAttachmentsByLessonID)

					// Upload endpoint - only admin and teacher can upload
					uploadGroup := attachmentGroup.Group("")
					uploadGroup.Use(h.requireAdminOrTeacher())
					{
						uploadGroup.POST("", h.UploadFile)
					}

					// Download endpoint - any authenticated user with access to the lesson can download
					attachmentGroup.GET("/:attachmentId", h.DownloadFile)

					// Delete attachment - only admin and teacher can delete
					deleteGroup := attachmentGroup.Group("/:attachmentId")
					deleteGroup.Use(h.requireAdminOrTeacher())
					{
						deleteGroup.DELETE("", h.DeleteAttachment)
					}
				}
			}
		}
	}

	// Keep the old routes for backward compatibility
	oldAttachmentGroup := router.Group("/api/v1/attachments")
	oldAttachmentGroup.Use(middleware.AuthMiddleware(h.authService))
	{
		// Upload endpoint - only admin and teacher can upload
		uploadGroup := oldAttachmentGroup.Group("/upload")
		uploadGroup.Use(h.requireAdminOrTeacher())
		{
			uploadGroup.POST("/:lessonId", h.UploadFile)
		}

		// Download endpoint - any authenticated user with access to the lesson can download
		oldAttachmentGroup.GET("/download/:id", h.DownloadFile)

		// Get attachments for a lesson
		oldAttachmentGroup.GET("/lesson/:lessonId", h.GetAttachmentsByLessonID)

		// Delete attachment - only admin and teacher can delete
		deleteGroup := oldAttachmentGroup.Group("/delete")
		deleteGroup.Use(h.requireAdminOrTeacher())
		{
			deleteGroup.DELETE("/:id", h.DeleteAttachment)
		}
	}
}

// requireAdminOrTeacher is a middleware that checks if the user has admin or teacher
func (h *AttachmentHandler) requireAdminOrTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the claims from the context
		claims, exists := c.Get("claims")
		if !exists {
			middleware.RespondWithError(c, http.StatusUnauthorized, "Authentication required")
			c.Abort()
			return
		}

		// Check if the user has the required role
		keycloakClaims, ok := claims.(*services.KeycloakClaims)
		if !ok {
			middleware.RespondWithError(c, http.StatusInternalServerError, "Invalid claims type")
			c.Abort()
			return
		}

		if !h.authService.HasRole(keycloakClaims, "admin") && !h.authService.HasRole(keycloakClaims, "teacher") {
			middleware.RespondWithError(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// UploadFile handles POST /api/v1/courses/:id/chapters/:chapterId/lessons/:lessonId/attachments
// @Summary Upload a file
// @Description Upload a file to a lesson
// @Tags attachments
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param chapterId path int true "Chapter ID"
// @Param lessonId path int true "Lesson ID"
// @Param file formData file true "File to upload"
// @Success 201 {object} map[string]interface{} "File uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid course, chapter, or lesson ID or file"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/chapters/{chapterId}/lessons/{lessonId}/attachments [post]
func (h *AttachmentHandler) UploadFile(c *gin.Context) {
	// Parse course ID, chapter ID, and lesson ID
	courseIdStr := c.Param("id")
	chapterIdStr := c.Param("chapterId")
	lessonIdStr := c.Param("lessonId")

	// For backward compatibility, if courseId or chapterId is not provided, use the lessonId directly
	if courseIdStr == "" || chapterIdStr == "" {
		lessonIdStr = c.Param("lessonId")
		if lessonIdStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid lesson ID",
			})
			return
		}

		lessonId, err := strconv.ParseUint(lessonIdStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid lesson ID",
			})
			return
		}

		// Get the file from the request
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "No file uploaded",
			})
			return
		}

		// Upload the file
		uploadResponse, err := h.service.UploadFile(file, 0, 0, uint(lessonId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"error":   false,
			"data":    uploadResponse,
			"message": "File uploaded successfully",
		})
		return
	}

	// Parse IDs
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
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

	lessonId, err := strconv.ParseUint(lessonIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid lesson ID",
		})
		return
	}

	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "No file uploaded",
		})
		return
	}

	// Upload the file
	uploadResponse, err := h.service.UploadFile(file, uint(courseId), uint(chapterId), uint(lessonId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"data":    uploadResponse,
		"message": "File uploaded successfully",
	})
}

// DownloadFile handles GET /api/v1/courses/:id/chapters/:chapterId/lessons/:lessonId/attachments/:attachmentId
// @Summary Download a file
// @Description Download a file by its ID
// @Tags attachments
// @Produce octet-stream
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param chapterId path int true "Chapter ID"
// @Param lessonId path int true "Lesson ID"
// @Param attachmentId path int true "Attachment ID"
// @Success 200 {file} binary "File content"
// @Failure 400 {object} map[string]interface{} "Invalid attachment ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "File not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/chapters/{chapterId}/lessons/{lessonId}/attachments/{attachmentId} [get]
func (h *AttachmentHandler) DownloadFile(c *gin.Context) {
	// Parse attachment ID
	idStr := c.Param("attachmentId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid attachment ID",
		})
		return
	}

	// Get the attachment and MinIO object
	attachment, object, err := h.service.DownloadFile(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "attachment not found" || err.Error() == "file not found in MinIO" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	defer object.Close()

	// Get the user ID from the context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "User not authenticated",
		})
		return
	}

	// Check if the user has access to the lesson
	hasAccess, err := h.service.HasAccessToLesson(userID.(uint), attachment.LessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You don't have access to this lesson",
		})
		return
	}

	// Set the appropriate headers
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachment.Name))
	c.Header("Content-Type", "application/octet-stream")

	// Copy the object to the response writer
	if _, err := io.Copy(c.Writer, object); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to stream file: " + err.Error(),
		})
		return
	}
}

// GetAttachmentsByLessonID handles GET /api/v1/courses/:id/chapters/:chapterId/lessons/:lessonId/attachments
// @Summary Get attachments for a lesson
// @Description Get all attachments for a specific lesson
// @Tags attachments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param chapterId path int true "Chapter ID"
// @Param lessonId path int true "Lesson ID"
// @Success 200 {object} map[string]interface{} "Returns a list of attachments"
// @Failure 400 {object} map[string]interface{} "Invalid course, chapter, or lesson ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/chapters/{chapterId}/lessons/{lessonId}/attachments [get]
func (h *AttachmentHandler) GetAttachmentsByLessonID(c *gin.Context) {
	// Parse course ID, chapter ID, and lesson ID
	courseIdStr := c.Param("id")
	chapterIdStr := c.Param("chapterId")
	lessonIdStr := c.Param("lessonId")

	// For backward compatibility, if courseId or chapterId is not provided, use the lessonId directly
	if courseIdStr == "" || chapterIdStr == "" {
		lessonIdStr = c.Param("lessonId")
		if lessonIdStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid lesson ID",
			})
			return
		}

		lessonId, err := strconv.ParseUint(lessonIdStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid lesson ID",
			})
			return
		}

		// Get the attachments
		attachments, err := h.service.GetAttachmentsByLessonID(0, 0, uint(lessonId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"error": false,
			"data":  attachments,
		})
		return
	}

	// Parse IDs
	courseId, err := strconv.ParseUint(courseIdStr, 10, 32)
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

	lessonId, err := strconv.ParseUint(lessonIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid lesson ID",
		})
		return
	}

	// Get the attachments
	attachments, err := h.service.GetAttachmentsByLessonID(uint(courseId), uint(chapterId), uint(lessonId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  attachments,
	})
}

// DeleteAttachment handles DELETE /api/v1/courses/:id/chapters/:chapterId/lessons/:lessonId/attachments/:attachmentId
// @Summary Delete an attachment
// @Description Delete an attachment by its ID
// @Tags attachments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Param chapterId path int true "Chapter ID"
// @Param lessonId path int true "Lesson ID"
// @Param attachmentId path int true "Attachment ID"
// @Success 200 {object} map[string]interface{} "Attachment deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid attachment ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Attachment not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /courses/{id}/chapters/{chapterId}/lessons/{lessonId}/attachments/{attachmentId} [delete]
func (h *AttachmentHandler) DeleteAttachment(c *gin.Context) {
	// Parse attachment ID
	idStr := c.Param("attachmentId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid attachment ID",
		})
		return
	}

	// Delete the attachment
	err = h.service.DeleteAttachment(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "attachment not found" {
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
		"message": "Attachment deleted successfully",
	})
}
