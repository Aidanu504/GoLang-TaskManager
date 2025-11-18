package utils

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// General error function
func Error(c *gin.Context, status int, message string) {
    c.JSON(status, gin.H{"error": message})
}

// helper function for 500 errors
func ServerError(c *gin.Context, message string) {
    Error(c, http.StatusInternalServerError, message)
}