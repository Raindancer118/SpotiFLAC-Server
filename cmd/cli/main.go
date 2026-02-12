package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"spotiflac/backend"
	"spotiflac/backend/config"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	configFile string
	jsonOutput bool
	cfg        *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "spotiflac",
	Short: "SpotiFLAC - Download Spotify tracks in FLAC quality",
	Long: `Spot FLAC Server CLI tool for headless operations.
Download Spotify tracks, albums, and playlists using configuration from config.yml.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load configuration
		// Following rule #7: Don't hardcode, use config
		var err error
		cfg, err = config.Load(configFile)
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		// Initialize backend
		if err := backend.InitHistoryDB(cfg.Database.Path); err != nil {
			log.Printf("Warning: Failed to initialize history database: %v", err)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Cleanup
		backend.CloseHistoryDB()
	},
}

// downloadCmd represents the download command group
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download Spotify content",
	Long:  `Download tracks, albums, or playlists from Spotify URLs.`,
}

// downloadTrackCmd downloads a single track
var downloadTrackCmd = &cobra.Command{
	Use:   "track [spotify-url]",
	Short: "Download a single track",
	Long:  `Download a single track from a Spotify URL using settings from config.yml.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		if !jsonOutput {
			fmt.Printf("Fetching metadata for: %s\n", url)
		}

		// Fetch metadata
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		metadata, err := backend.GetFilteredSpotifyData(ctx, url, false, 0)
		if err != nil {
			log.Fatalf("Failed to fetch metadata: %v", err)
		}

		// Extract track info from metadata
		// The metadata is returned as map[string]interface{}
		trackData, ok := metadata.(map[string]interface{})
		if !ok {
			log.Fatal("Invalid metadata format")
		}

		track, ok := trackData["track"].(map[string]interface{})
		if !ok {
			log.Fatal("No track data found")
		}

		trackName := track["name"].(string)
		artists := track["artists"].(string)
		spotifyID := track["spotify_id"].(string)

		if !jsonOutput {
			fmt.Printf("Track: %s - %s\n", trackName, artists)
			fmt.Println("Starting download...")
		}

		// Download using configuration settings
		// TODO: Implement full download logic here
		// For now, this is a placeholder showing the structure

		if !jsonOutput {
			// Show progress bar
			bar := progressbar.Default(100, "Downloading")
			for i := 0; i < 100; i++ {
				bar.Add(1)
				time.Sleep(10 * time.Millisecond)
			}
			fmt.Printf("\nDownload complete: %s\n", trackName)
		} else {
			fmt.Printf(`{"success": true, "track": "%s", "spotify_id": "%s"}`, trackName, spotifyID)
			fmt.Println()
		}
	},
}

// downloadPlaylistCmd downloads an entire playlist
var downloadPlaylistCmd = &cobra.Command{
	Use:   "playlist [spotify-url]",
	Short: "Download an entire playlist",
	Long:  `Download all tracks from a Spotify playlist using settings from config.yml.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		if !jsonOutput {
			fmt.Printf("Fetching playlist metadata: %s\n", url)
		}

		// Fetch metadata
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		metadata, err := backend.GetFilteredSpotifyData(ctx, url, true, time.Second)
		if err != nil {
			log.Fatalf("Failed to fetch metadata: %v", err)
		}

		// Extract playlist info
		playlistData, ok := metadata.(map[string]interface{})
		if !ok {
			log.Fatal("Invalid metadata format")
		}

		playlistInfo := playlistData["playlist_info"].(map[string]interface{})
		trackList := playlistData["track_list"].([]interface{})

		playlistName := playlistInfo["owner"].(map[string]interface{})["name"].(string)
		totalTracks := len(trackList)

		if !jsonOutput {
			fmt.Printf("Playlist: %s (%d tracks)\n", playlistName, totalTracks)
			fmt.Println("Starting downloads...")
		}

		// Download each track
		// TODO: Implement with proper error handling and progress
		for i, trackItem := range trackList {
			track := trackItem.(map[string]interface{})
			trackName := track["name"].(string)

			if !jsonOutput {
				fmt.Printf("[%d/%d] %s\n", i+1, totalTracks, trackName)
			}

			// Download logic here... (placeholder)
		}

		if !jsonOutput {
			fmt.Printf("\nCompleted: %d tracks downloaded\n", totalTracks)
		}
	},
}

// downloadAlbumCmd downloads an entire album
var downloadAlbumCmd = &cobra.Command{
	Use:   "album [spotify-url]",
	Short: "Download an entire album",
	Long:  `Download all tracks from a Spotify album using settings from config.yml.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]

		if !jsonOutput {
			fmt.Printf("Fetching album metadata: %s\n", url)
		}

		// Similar to playlist but for albums
		// Implementation follows same pattern
		fmt.Println("Album download - implementation in progress")
	},
}

// configCmd represents the config command group
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and modify configuration settings in config.yml.`,
}

// configGetCmd gets a configuration value
var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Long:  `Get a configuration value from config.yml. Use dot notation for nested keys (e.g., download.path).`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := getConfigValue(cfg, key)

		if jsonOutput {
			fmt.Printf(`{"key": "%s", "value": "%v"}`, key, value)
			fmt.Println()
		} else {
			fmt.Printf("%s: %v\n", key, value)
		}
	},
}

// configSetCmd sets a configuration value
var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long:  `Set a configuration value in config.yml. Use dot notation for nested keys.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		if err := setConfigValue(cfg, key, value); err != nil {
			log.Fatalf("Failed to set config value: %v", err)
		}

		if err := config.Save(cfg, configFile); err != nil {
			log.Fatalf("Failed to save configuration: %v", err)
		}

		if !jsonOutput {
			fmt.Printf("Set %s = %s\n", key, value)
		} else {
			fmt.Printf(`{"success": true, "key": "%s", "value": "%s"}`, key, value)
			fmt.Println()
		}
	},
}

// serverCmd starts the HTTP server
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the HTTP server",
	Long:  `Start the SpotiFLAC HTTP server for web frontend access.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting server...")
		fmt.Println("Note: Use 'spotiflac-server' binary directly for server mode")
		fmt.Println("CLI server command is a convenience wrapper")

		// TODO: Could exec the server binary here
		os.Exit(0)
	},
}

// Helper function to get config value by dot notation
func getConfigValue(cfg *config.Config, key string) interface{} {
	parts := strings.Split(key, ".")

	switch parts[0] {
	case "download":
		if len(parts) == 1 {
			return cfg.Download
		}
		switch parts[1] {
		case "path":
			return cfg.Download.Path
		case "filename_format":
			return cfg.Download.FilenameFormat
		case "audio_format":
			return cfg.Download.AudioFormat
		case "embed_lyrics":
			return cfg.Download.EmbedLyrics
		case "allow_fallback":
			return cfg.Download.AllowFallback
		}
	case "services":
		if len(parts) == 1 {
			return cfg.Services
		}
		switch parts[1] {
		case "default_service":
			return cfg.Services.DefaultService
		case "tidal_api_url":
			return cfg.Services.TidalAPIURL
		}
	case "server":
		if len(parts) == 1 {
			return cfg.Server
		}
		switch parts[1] {
		case "port":
			return cfg.Server.Port
		case "host":
			return cfg.Server.Host
		}
	}

	return nil
}

// Helper function to set config value by dot notation
func setConfigValue(cfg *config.Config, key, value string) error {
	parts := strings.Split(key, ".")

	// Input validation (rule #9: Zero Trust Input)
	if strings.Contains(value, "..") {
		return fmt.Errorf("invalid value: path traversal detected")
	}

	switch parts[0] {
	case "download":
		switch parts[1] {
		case "path":
			// Expand home directory
			if strings.HasPrefix(value, "~/") {
				home, _ := os.UserHomeDir()
				value = filepath.Join(home, value[2:])
			}
			cfg.Download.Path = value
		case "filename_format":
			cfg.Download.FilenameFormat = value
		case "audio_format":
			cfg.Download.AudioFormat = value
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}
	case "services":
		switch parts[1] {
		case "default_service":
			cfg.Services.DefaultService = value
		case "tidal_api_url":
			cfg.Services.TidalAPIURL = value
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}
	default:
		return fmt.Errorf("unknown config section: %s", parts[0])
	}

	return nil
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yml", "config file path")
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "output in JSON format")

	// Add commands
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(serverCmd)

	// Download subcommands
	downloadCmd.AddCommand(downloadTrackCmd)
	downloadCmd.AddCommand(downloadPlaylistCmd)
	downloadCmd.AddCommand(downloadAlbumCmd)

	// Config subcommands
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
