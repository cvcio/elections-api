package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response : JSON Response Object
type Response struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code"`
}

// ResponseJSON ...
func ResponseJSON(c *gin.Context, d interface{}) {
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    d,
		Code:    http.StatusOK,
	})
}

// ResponseError ...
func ResponseError(c *gin.Context, s int, e string) {
	c.JSON(s, &Response{
		Success: false,
		Error:   e,
		Code:    s,
	})
	c.Abort()
}
