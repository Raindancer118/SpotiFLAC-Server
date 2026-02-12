package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spotiflac/backend"
	"spotiflac/backend/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler wraps the application logic and exposes it via HTTP
// All methods from app.go are available as HTTP endpoints here
type Handler struct {
	// We could embed the App struct here, but we'll access backend directly
	// to avoid Wails dependencies
}

// NewHandler creates a new API handler
func NewHandler() *Handler {
	return &Handler{}
}

// HealthCheck returns server health status
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().Unix(),
	})
}

// GetSpotifyMetadata handles metadata fetching from Spotify URLs
// Endpoint: POST /api/spotify/metadata
func (h *Handler) GetSpotifyMetadata(c *gin.Context) {
	var req struct {
		URL     string  `json:"url" binding:"required"`
		Batch   bool    `json:"batch"`
		Delay   float64 `json:"delay"`
		Timeout float64 `json:"timeout"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Input validation (rule #9: Zero Trust Input)
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
		return
	}

	// Sanitize URL to prevent injection attacks
	req.URL = strings.TrimSpace(req.URL)

	if req.Delay == 0 {
		req.Delay = 1.0
	}
	if req.Timeout == 0 {
		req.Timeout = 300.0
	}

	// Get configuration
	cfg := config.Get()

	// Create context with timeout
	ctx := c.Request.Context()

	// Check if we should use SpotFetch API
	var data interface{}
	var err error

	if cfg.Services.UseSpotFetchAPI && cfg.Services.SpotFetchAPIURL != "" {
		data, err = backend.GetSpotifyDataWithAPI(
			ctx,
			req.URL,
			true,
			cfg.Services.SpotFetchAPIURL,
			req.Batch,
			time.Duration(req.Delay*float64(time.Second)),
		)
	} else {
		data, err = backend.GetFilteredSpotifyData(
			ctx,
			req.URL,
			req.Batch,
			time.Duration(req.Delay*float64(time.Second)),
		)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch metadata: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}

// SearchSpotify handles Spotify search requests
// Endpoint: POST /api/spotify/search
func (h *Handler) SearchSpotify(c *gin.Context) {
	var req struct {
		Query string `json:"query" binding:"required"`
		Limit int    `json:"limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Input validation (rule #9)
	if req.Query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Sanitize query to prevent SQL injection and XSS
	req.Query = strings.TrimSpace(req.Query)

	if req.Limit <= 0 {
		req.Limit = 10
	}

	ctx := c.Request.Context()
	result, err := backend.SearchSpotify(ctx, req.Query, req.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Search failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// SearchSpotifyByType handles type-specific Spotify searches
// Endpoint: POST /api/spotify/search-by-type
func (h *Handler) SearchSpotifyByType(c *gin.Context) {
	var req struct {
		Query      string `json:"query" binding:"required"`
		SearchType string `json:"search_type" binding:"required"`
		Limit      int    `json:"limit"`
		Offset     int    `json:"offset"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Input validation
	if req.Query == "" || req.SearchType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query and search type are required"})
		return
	}

	if req.Limit <= 0 {
		req.Limit = 50
	}

	ctx := c.Request.Context()
	results, err := backend.SearchSpotifyByType(ctx, req.Query, req.SearchType, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Search failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetStreamingURLs retrieves streaming URLs for a track
// Endpoint: POST /api/spotify/streaming-urls
func (h *Handler) GetStreamingURLs(c *gin.Context) {
	var req struct {
		SpotifyTrackID string `json:"spotify_track_id" binding:"required"`
		Region         string `json:"region"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if req.SpotifyTrackID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Spotify track ID is required"})
		return
	}

	if req.Region == "" {
		req.Region = "US"
	}

	client := backend.NewSongLinkClient()
	urls, err := client.GetAllURLsFromSpotify(req.SpotifyTrackID, req.Region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get streaming URLs: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, urls)
}

// GetSettings returns current configuration settings
// Endpoint: GET /api/settings
func (h *Handler) GetSettings(c *gin.Context) {
	cfg := config.Get()

	// Convert config to map format expected by frontend
	settings := map[string]interface{}{
		"downloadPath":         cfg.Download.Path,
		"filenameFormat":       cfg.Download.FilenameFormat,
		"audioFormat":          cfg.Download.AudioFormat,
		"embedLyrics":          cfg.Download.EmbedLyrics,
		"embedMaxQualityCover": cfg.Download.EmbedMaxQualityCover,
		"trackNumber":          cfg.Download.TrackNumber,
		"useAlbumTrackNumber":  cfg.Download.UseAlbumTrackNumber,
		"useFirstArtistOnly":   cfg.Download.UseFirstArtistOnly,
		"allowFallback":        cfg.Download.AllowFallback,
		"defaultService":       cfg.Services.DefaultService,
		"tidalAPIUrl":          cfg.Services.TidalAPIURL,
		"useSpotFetchAPI":      cfg.Services.UseSpotFetchAPI,
		"spotFetchAPIUrl":      cfg.Services.SpotFetchAPIURL,
		"theme":                cfg.UI.Theme,
		"themeMode":            cfg.UI.ThemeMode,
		"fontFamily":           cfg.UI.FontFamily,
	}

	c.JSON(http.StatusOK, settings)
}

// SaveSettings updates configuration settings
// Endpoint: POST /api/settings
func (h *Handler) SaveSettings(c *gin.Context) {
	var settings map[string]interface{}

	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid settings format"})
		return
	}

	// Get current config
	cfg := config.Get()

	// Update fields that are present in the request
	// Following rule #9: validate all input before applying
	if val, ok := settings["downloadPath"].(string); ok {
		// Validate path (prevent path traversal)
		if strings.Contains(val, "..") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid download path"})
			return
		}
		cfg.Download.Path = val
	}

	if val, ok := settings["filenameFormat"].(string); ok {
		cfg.Download.FilenameFormat = val
	}

	if val, ok := settings["audioFormat"].(string); ok {
		cfg.Download.AudioFormat = val
	}

	if val, ok := settings["embedLyrics"].(bool); ok {
		cfg.Download.EmbedLyrics = val
	}

	if val, ok := settings["embedMaxQualityCover"].(bool); ok {
		cfg.Download.EmbedMaxQualityCover = val
	}

	if val, ok := settings["trackNumber"].(bool); ok {
		cfg.Download.TrackNumber = val
	}

	if val, ok := settings["useAlbumTrackNumber"].(bool); ok {
		cfg.Download.UseAlbumTrackNumber = val
	}

	if val, ok := settings["useFirstArtistOnly"].(bool); ok {
		cfg.Download.UseFirstArtistOnly = val
	}

	if val, ok := settings["allowFallback"].(bool); ok {
		cfg.Download.AllowFallback = val
	}

	if val, ok := settings["defaultService"].(string); ok {
		cfg.Services.DefaultService = val
	}

	if val, ok := settings["tidalAPIUrl"].(string); ok {
		cfg.Services.TidalAPIURL = val
	}

	if val, ok := settings["useSpotFetchAPI"].(bool); ok {
		cfg.Services.UseSpotFetchAPI = val
	}

	if val, ok := settings["spotFetchAPIUrl"].(string); ok {
		cfg.Services.SpotFetchAPIURL = val
	}

	if val, ok := settings["theme"].(string); ok {
		cfg.UI.Theme = val
	}

	if val, ok := settings["themeMode"].(string); ok {
		cfg.UI.ThemeMode = val
	}

	if val, ok := settings["fontFamily"].(string); ok {
		cfg.UI.FontFamily = val
	}

	// Save to config.yml
	if err := config.Save(cfg, "config.yml"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save settings",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetDefaults returns default system values
// Endpoint: GET /api/defaults
func (h *Handler) GetDefaults(c *gin.Context) {
	defaults := map[string]string{
		"downloadPath": backend.GetDefaultMusicPath(),
	}

	c.JSON(http.StatusOK, defaults)
}

// DownloadResponse matches the response structure from app.go
type DownloadResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	File          string `json:"file,omitempty"`
	Error         string `json:"error,omitempty"`
	AlreadyExists bool   `json:"already_exists,omitempty"`
	ItemID        string `json:"item_id,omitempty"`
}

// DownloadTrack handles track download requests
// Endpoint: POST /api/download/track
// This is a complex endpoint that will be continued in the next file section
func (h *Handler) DownloadTrack(c *gin.Context) {
	// This will be implemented similarly to app.go's DownloadTrack method
	// but adapted for HTTP context
	var req map[string]interface{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// For now, return a placeholder
	// Full implementation will follow the app.go logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Download endpoint - implementation in progress",
		"request": req,
	})
}

// GetDownloadQueue returns the current download queue status
// Endpoint: GET /api/download/queue
func (h *Handler) GetDownloadQueue(c *gin.Context) {
	queue := backend.GetDownloadQueue()
	c.JSON(http.StatusOK, queue)
}

// GetDownloadProgress returns download progress information
// Endpoint: GET /api/download/progress
func (h *Handler) GetDownloadProgress(c *gin.Context) {
	progress := backend.GetDownloadProgress()
	c.JSON(http.StatusOK, progress)
}

// ClearCompletedDownloads clears completed items from queue
// Endpoint: POST /api/download/queue/clear
func (h *Handler) ClearCompletedDownloads(c *gin.Context) {
	backend.ClearDownloadQueue()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// ClearAllDownloads clears all downloads from queue
// Endpoint: POST /api/download/queue/clear-all
func (h *Handler) ClearAllDownloads(c *gin.Context) {
	backend.ClearAllDownloads()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// CancelAllQueuedItems cancels all queued download items
// Endpoint: POST /api/download/queue/cancel-all
func (h *Handler) CancelAllQueuedItems(c *gin.Context) {
	backend.CancelAllQueuedItems()
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetDownloadHistory returns download history
// Endpoint: GET /api/history/downloads
func (h *Handler) GetDownloadHistory(c *gin.Context) {
	history, err := backend.GetHistoryItems("SpotiFLAC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get download history",
		})
		return
	}

	c.JSON(http.StatusOK, history)
}

// ClearDownloadHistory clears download history
// Endpoint: POST /api/history/downloads/clear
func (h *Handler) ClearDownloadHistory(c *gin.Context) {
	if err := backend.ClearHistory("SpotiFLAC"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to clear history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DeleteDownloadHistoryItem deletes a specific download history item
// Endpoint: DELETE /api/history/downloads/:id
func (h *Handler) DeleteDownloadHistoryItem(c *gin.Context) {
	id := c.Param("id")

	if err := backend.DeleteHistoryItem(id, "SpotiFLAC"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete history item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// CheckFFmpegInstalled checks if FFmpeg is installed
// Endpoint: GET /api/system/ffmpeg/status
func (h *Handler) CheckFFmpegInstalled(c *gin.Context) {
	installed, err := backend.IsFFmpegInstalled()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check FFmpeg status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"installed": installed,
	})
}

// AnalyzeTrack analyzes an audio file
// Endpoint: POST /api/analysis/track
func (h *Handler) AnalyzeTrack(c *gin.Context) {
	var req struct {
		FilePath string `json:"file_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Validate file path (rule #9: prevent path traversal)
	if strings.Contains(req.FilePath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file path"})
		return
	}

	result, err := backend.AnalyzeTrack(req.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to analyze track: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Helper function to convert interface{} to JSON string
func toJSONString(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
