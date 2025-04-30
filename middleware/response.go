package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Set("SuccessResponse", SuccessResponse)
		c.Set("ErrorResponse", ErrorResponse)
		c.Set("CreatedResponse", CreatedResponse)
		c.Set("NotFoundResponse", NotFoundResponse)
		c.Set("BadRequestResponse", BadRequestResponse)
		c.Set("InternalServerErrorResponse", InternalServerErrorResponse)

		c.Next()
	}
}

func SuccessResponse(c *gin.Context, data interface{}, message string) {
	response := Response{
		Error:   false,
		Data:    data,
		Message: message,
	}
	c.JSON(http.StatusOK, response)
}

func CreatedResponse(c *gin.Context, data interface{}, message string) {
	response := Response{
		Error:   false,
		Data:    data,
		Message: message,
	}
	c.JSON(http.StatusCreated, response)
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	response := Response{
		Error:   true,
		Message: message,
	}
	c.JSON(statusCode, response)
}

func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message)
}
