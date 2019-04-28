package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ContentTypeCheck(c *gin.Context) {
	if !strings.Contains(c.GetHeader("content-Type"), "multipart/form-data") {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"status": http.StatusUnsupportedMediaType, "result": "not support " + c.GetHeader("content-Type")})
		c.Abort()
	}
	c.Next()
}
