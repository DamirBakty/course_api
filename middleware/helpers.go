package middleware

import (
	"github.com/gin-gonic/gin"
)

func GetSuccessResponse(c *gin.Context) func(*gin.Context, interface{}, string) {
	return c.MustGet("SuccessResponse").(func(*gin.Context, interface{}, string))
}

func GetErrorResponse(c *gin.Context) func(*gin.Context, int, string) {
	return c.MustGet("ErrorResponse").(func(*gin.Context, int, string))
}

func GetCreatedResponse(c *gin.Context) func(*gin.Context, interface{}, string) {
	return c.MustGet("CreatedResponse").(func(*gin.Context, interface{}, string))
}

func GetNotFoundResponse(c *gin.Context) func(*gin.Context, string) {
	return c.MustGet("NotFoundResponse").(func(*gin.Context, string))
}

func GetBadRequestResponse(c *gin.Context) func(*gin.Context, string) {
	return c.MustGet("BadRequestResponse").(func(*gin.Context, string))
}

func GetInternalServerErrorResponse(c *gin.Context) func(*gin.Context, string) {
	return c.MustGet("InternalServerErrorResponse").(func(*gin.Context, string))
}

func RespondWithSuccess(c *gin.Context, data interface{}, message string) {
	GetSuccessResponse(c)(c, data, message)
}

func RespondWithError(c *gin.Context, statusCode int, message string) {
	GetErrorResponse(c)(c, statusCode, message)
}

func RespondWithCreated(c *gin.Context, data interface{}, message string) {
	GetCreatedResponse(c)(c, data, message)
}

func RespondWithNotFound(c *gin.Context, message string) {
	GetNotFoundResponse(c)(c, message)
}

func RespondWithBadRequest(c *gin.Context, message string) {
	GetBadRequestResponse(c)(c, message)
}

func RespondWithInternalServerError(c *gin.Context, message string) {
	GetInternalServerErrorResponse(c)(c, message)
}
