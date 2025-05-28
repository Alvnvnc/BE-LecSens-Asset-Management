package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs HTTP requests and responses
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Save request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create a response writer that captures the response
		responseWriter := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Extract user and tenant information from context
		userID, userExists := c.Get("user_id")
		tenantID, tenantExists := c.Get("tenant_id")

		// Log request details
		log.Printf(
			"Request: %s %s | Status: %d | Duration: %v | User: %v | Tenant: %v | IP: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			userID,
			tenantID,
			c.ClientIP(),
		)

		// Log extra details for non-success responses
		if c.Writer.Status() >= 400 {
			log.Printf(
				"Error Details - Method: %s | Path: %s | Status: %d | User: %v | Tenant: %v | Response: %s",
				c.Request.Method,
				c.Request.URL.Path,
				c.Writer.Status(),
				userExists,
				tenantExists,
				responseWriter.body.String(),
			)
		}
	}
}

// responseBodyWriter is a wrapper around gin.ResponseWriter that captures the response body
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body
func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString captures the response body
func (w *responseBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
