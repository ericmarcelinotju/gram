package response

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SetResponse struct for response
type SetResponse struct {
	Status     string      `json:"status"`
	Data       interface{} `json:"data,omitempty"`
	Code       int         `json:"code"`
	AccessTime string      `json:"accessTime"`
}

// ResponseSuccess for endpoint success
func ResponseSuccess(c *gin.Context, data interface{}) {
	response := SetResponse{
		Status:     http.StatusText(http.StatusOK),
		AccessTime: time.Now().Format("02-01-2006 15:04:05"),
		Data:       data,
		Code:       http.StatusOK,
	}

	c.JSON(http.StatusOK, response)
}

// ResponseError for endpoint error
func ResponseError(c *gin.Context, err error, code int) {
	response := SetResponse{
		Status:     http.StatusText(code),
		AccessTime: time.Now().Format("02-01-2006 15:04:05"),
		Data:       err.Error(),
		Code:       code,
	}
	c.JSON(code, response)
}

// ResponseError for endpoint error
func ResponseAbort(c *gin.Context, err error, code int) {
	response := SetResponse{
		Status:     http.StatusText(code),
		AccessTime: time.Now().Format("02-01-2006 15:04:05"),
		Data:       err.Error(),
		Code:       code,
	}
	c.AbortWithStatusJSON(code, response)
}

// ResponseImage endpoint image
func ResponseFile(c *gin.Context, filename string, file []byte) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", len(file)))

	if _, err := c.Writer.Write(file); err != nil {
		ResponseError(c, err, http.StatusNotFound)
		return
	}
}

// ResponseHTML for HTML endpoint
func ResponseHTML(c *gin.Context, html string) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
