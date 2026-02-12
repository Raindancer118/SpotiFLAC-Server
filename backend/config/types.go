package config

// Config represents the complete application configuration
// loaded from config.yml
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Download DownloadConfig `yaml:"download"`
	Services ServicesConfig `yaml:"services"`
	UI       UIConfig       `yaml:"ui"`
	Database DatabaseConfig `yaml:"database"`
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Host        string   `yaml:"host"`
	Port        int      `yaml:"port"`
	CORSOrigins []string `yaml:"cors_origins"`
}

// DownloadConfig contains download preferences
// These settings match the settings used in the Wails frontend
type DownloadConfig struct {
	Path                  string `yaml:"path"`
	FilenameFormat        string `yaml:"filename_format"`
	AudioFormat           string `yaml:"audio_format"`
	EmbedLyrics           bool   `yaml:"embed_lyrics"`
	EmbedMaxQualityCover  bool   `yaml:"embed_max_quality_cover"`
	TrackNumber           bool   `yaml:"track_number"`
	UseAlbumTrackNumber   bool   `yaml:"use_album_track_number"`
	UseFirstArtistOnly    bool   `yaml:"use_first_artist_only"`
	AllowFallback         bool   `yaml:"allow_fallback"`
}

// ServicesConfig contains streaming service settings
type ServicesConfig struct {
	DefaultService    string `yaml:"default_service"`
	TidalAPIURL       string `yaml:"tidal_api_url"`
	UseSpotFetchAPI   bool   `yaml:"use_spotfetch_api"`
	SpotFetchAPIURL   string `yaml:"spotfetch_api_url"`
}

// UIConfig contains user interface preferences
type UIConfig struct {
	Theme      string `yaml:"theme"`
	ThemeMode  string `yaml:"theme_mode"`
	FontFamily string `yaml:"font_family"`
}

// DatabaseConfig contains database settings
type DatabaseConfig struct {
	Path string `yaml:"path"`
}
