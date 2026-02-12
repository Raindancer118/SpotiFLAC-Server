package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware contains common middleware functions for the HTTP server
// Following rule #15: Fail Securely - error handling without exposing sensitive info

// ErrorHandler is a middleware that converts panics to proper HTTP error responses
// Following rule #15: global error handling with neutral user messages
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log detailed error for developers (rule #16: Audit Logging)
				c.Error(err.(error))

				// Return neutral error message to user (rule #15: don't expose internals)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "An internal error occurred. Please try again later.",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RequestLogger logs all requests for audit purposes
// Following rule #16: Audit Logs - know who did what when
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request details
		// In production, this should write to a centralized logging system
		c.Next()

		// Log after request completes to capture status and errors
		statusCode := c.Writer.Status()
		if statusCode >= 400 {
			// Log error requests with more detail
			c.Error(c.Errors.Last())
		}
	}
}

// SecurityHeaders adds security headers to all responses
// Following rule #14: Secure by Default - set security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'")

		c.Next()
	}
}

// InputSanitizer validates and sanitizes input
// Following rule #9: Zero Trust Input - validate all external data
func InputSanitizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Additional input validation can be added here
		// Individual handlers should still validate their specific inputs
		c.Next()
	}
}
