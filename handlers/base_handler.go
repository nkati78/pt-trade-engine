package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HTTPError struct {
	Message string `json:"message"`
}

type HTTPStatusCode int

const (
	HTTPStatusOK                  HTTPStatusCode = http.StatusOK
	HTTPStatusCreated             HTTPStatusCode = http.StatusCreated
	HTTPStatusBadRequest          HTTPStatusCode = http.StatusBadRequest
	HTTPStatusUnauthorized        HTTPStatusCode = http.StatusUnauthorized
	HTTPStatusInternalServerError HTTPStatusCode = http.StatusInternalServerError
)

// BaseHandler is a function that allows a handler to return a status code and a response
type BaseHandler func(c *gin.Context) (HTTPStatusCode, interface{})

// ToHandler converts a BaseHandler to a gin.HandlerFunc and allow
// it to be used in a gin router
func ToHandler(handler BaseHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		status, response := handler(c)
		c.JSON(int(status), response)
	}
}
