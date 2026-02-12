package server

import (
	"fmt"
	"log"
	"spotiflac/backend"
	"spotiflac/backend/config"
	"spotiflac/server/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	config *config.Config
}

// NewServer creates a new HTTP server instance
func NewServer(cfg *config.Config) *Server {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Add middleware
	router.Use(gin.Recovery())        // Panic recovery
	router.Use(api.RequestLogger())   // Request logging (rule #16)
	router.Use(api.ErrorHandler())    // Error handling (rule #15)
	router.Use(api.SecurityHeaders()) // Security headers (rule #14)
	router.Use(api.InputSanitizer())  // Input sanitization (rule #9)

	// Configure CORS (rule #13: Defense in Depth)
	corsConfig := cors.Config{
		AllowOrigins:     cfg.Server.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	router.Use(cors.New(corsConfig))

	return &Server{
		router: router,
		config: cfg,
	}
}

// SetupRoutes configures all HTTP routes
func (s *Server) SetupRoutes() {
	handler := api.NewHandler()

	// Health check
	s.router.GET("/health", handler.HealthCheck)

	// API routes
	apiGroup := s.router.Group("/api")
	{
		// Spotify metadata and search
		spotify := apiGroup.Group("/spotify")
		{
			spotify.POST("/metadata", handler.GetSpotifyMetadata)
			spotify.POST("/search", handler.SearchSpotify)
			spotify.POST("/search-by-type", handler.SearchSpotifyByType)
			spotify.POST("/streaming-urls", handler.GetStreamingURLs)
		}

		// Download operations
		download := apiGroup.Group("/download")
		{
			download.POST("/track", handler.DownloadTrack)
			download.GET("/queue", handler.GetDownloadQueue)
			download.GET("/progress", handler.GetDownloadProgress)
			download.POST("/queue/clear", handler.ClearCompletedDownloads)
			download.POST("/queue/clear-all", handler.ClearAllDownloads)
			download.POST("/queue/cancel-all", handler.CancelAllQueuedItems)
		}

		// History
		history := apiGroup.Group("/history")
		{
			history.GET("/downloads", handler.GetDownloadHistory)
			history.POST("/downloads/clear", handler.ClearDownloadHistory)
			history.DELETE("/downloads/:id", handler.DeleteDownloadHistoryItem)
		}

		// Settings
		apiGroup.GET("/settings", handler.GetSettings)
		apiGroup.POST("/settings", handler.SaveSettings)
		apiGroup.GET("/defaults", handler.GetDefaults)

		// System
		system := apiGroup.Group("/system")
		{
			system.GET("/ffmpeg/status", handler.CheckFFmpegInstalled)
		}

		// Analysis
		analysis := apiGroup.Group("/analysis")
		{
			analysis.POST("/track", handler.AnalyzeTrack)
		}
	}

	// WebSocket endpoint for real-time updates
	s.router.GET("/ws", api.HandleWebSocket)
}

// Start initializes and starts the HTTP server
func (s *Server) Start() error {
	// Initialize backend components
	if err := backend.InitHistoryDB(s.config.Database.Path); err != nil {
		return fmt.Errorf("failed to initialize history database: %w", err)
	}

	// Initialize WebSocket manager
	api.InitWebSocketManager()

	// Setup routes
	s.SetupRoutes()

	// Start server
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	log.Printf("Starting SpotiFLAC server on %s", addr)

	return s.router.Run(addr)
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	// Close database connections
	backend.CloseHistoryDB()

	log.Println("Server stopped")
	return nil
}
