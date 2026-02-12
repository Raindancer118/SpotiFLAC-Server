// HTTP API Client for SpotiFLAC Server
// Replaces Wails bindings with HTTP calls

// API Base URL - configurable via environment variable
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

/**
 * API Client class
 * Provides methods matching the original Wails bindings
 */
class ApiClient {
    private baseUrl: string;

    constructor(baseUrl: string = API_BASE_URL) {
        this.baseUrl = baseUrl;
    }

    /**
     * Generic fetch wrapper with error handling
     */
    private async fetch<T>(
        endpoint: string,
        options: RequestInit = {}
    ): Promise<T> {
        const url = `${this.baseUrl}${endpoint}`;

        const response = await fetch(url, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
        });

        if (!response.ok) {
            const error = await response.json().catch(() => ({ error: 'Request failed' }));
            throw new Error(error.error || `HTTP ${response.status}`);
        }

        return response.json();
    }

    // ==================== Spotify Metadata ====================

    async GetSpotifyMetadata(req: {
        url: string;
        batch?: boolean;
        delay?: number;
        timeout?: number;
    }): Promise<string> {
        const data = await this.fetch<any>('/api/spotify/metadata', {
            method: 'POST',
            body: JSON.stringify(req),
        });
        return JSON.stringify(data);
    }

    async SearchSpotify(req: {
        query: string;
        limit?: number;
    }): Promise<any> {
        return this.fetch('/api/spotify/search', {
            method: 'POST',
            body: JSON.stringify(req),
        });
    }

    async SearchSpotifyByType(req: {
        query: string;
        search_type: string;
        limit?: number;
        offset?: number;
    }): Promise<any[]> {
        return this.fetch('/api/spotify/search-by-type', {
            method: 'POST',
            body: JSON.stringify(req),
        });
    }

    async GetStreamingURLs(spotifyTrackID: string, region: string): Promise<string> {
        const data = await this.fetch<any>('/api/spotify/streaming-urls', {
            method: 'POST',
            body: JSON.stringify({ spotify_track_id: spotifyTrackID, region }),
        });
        return JSON.stringify(data);
    }

    // ==================== Downloads ====================

    async DownloadTrack(req: any): Promise<any> {
        return this.fetch('/api/download/track', {
            method: 'POST',
            body: JSON.stringify(req),
        });
    }

    async GetDownloadQueue(): Promise<any> {
        return this.fetch('/api/download/queue');
    }

    async GetDownloadProgress(): Promise<any> {
        return this.fetch('/api/download/progress');
    }

    async ClearCompletedDownloads(): Promise<void> {
        await this.fetch('/api/download/queue/clear', { method: 'POST' });
    }

    async ClearAllDownloads(): Promise<void> {
        await this.fetch('/api/download/queue/clear-all', { method: 'POST' });
    }

    async CancelAllQueuedItems(): Promise<void> {
        await this.fetch('/api/download/queue/cancel-all', { method: 'POST' });
    }

    async AddToDownloadQueue(
        spotifyID: string,
        trackName: string,
        artistName: string,
        albumName: string
    ): Promise<string> {
        // This is handled server-side, return a placeholder ID
        return `${spotifyID}-${Date.now()}`;
    }

    async MarkDownloadItemFailed(itemID: string, errorMsg: string): Promise<void> {
        // Would need a server endpoint for this
        console.warn('MarkDownloadItemFailed not implemented in HTTP mode');
    }

    async ExportFailedDownloads(): Promise<string> {
        // This would need server implementation or client-side handling
        console.warn('ExportFailedDownloads not implemented in HTTP mode');
        return 'Not implemented';
    }

    // ==================== History ====================

    async GetDownloadHistory(): Promise<any[]> {
        return this.fetch('/api/history/downloads');
    }

    async ClearDownloadHistory(): Promise<void> {
        await this.fetch('/api/history/downloads/clear', { method: 'POST' });
    }

    async DeleteDownloadHistoryItem(id: string): Promise<void> {
        await this.fetch(`/api/history/downloads/${id}`, { method: 'DELETE' });
    }

    async GetFetchHistory(): Promise<any[]> {
        // This is stored client-side in the current implementation
        // Keep using localStorage
        return [];
    }

    async AddFetchHistory(item: any): Promise<void> {
        // Client-side only
    }

    async ClearFetchHistory(): Promise<void> {
        // Client-side only
    }

    async DeleteFetchHistoryItem(id: string): Promise<void> {
        // Client-side only
    }

    async ClearFetchHistoryByType(itemType: string): Promise<void> {
        // Client-side only
    }

    // ==================== Settings ====================

    async LoadSettings(): Promise<Record<string, any>> {
        return this.fetch('/api/settings');
    }

    async SaveSettings(settings: Record<string, any>): Promise<void> {
        await this.fetch('/api/settings', {
            method: 'POST',
            body: JSON.stringify(settings),
        });
    }

    async GetDefaults(): Promise<Record<string, string>> {
        return this.fetch('/api/defaults');
    }

    // ==================== System ====================

    async CheckFFmpegInstalled(): Promise<boolean> {
        const data = await this.fetch<{ installed: boolean }>('/api/system/ffmpeg/status');
        return data.installed;
    }

    async DownloadFFmpeg(): Promise<{ success: boolean; error?: string }> {
        // FFmpeg installation should be done server-side
        console.warn('DownloadFFmpeg not available in HTTP mode - install on server');
        return { success: false, error: 'Not available in HTTP mode' };
    }

    async GetOSInfo(): Promise<string> {
        // Not needed in HTTP mode
        return 'Server';
    }

    // ==================== File Operations ====================

    async OpenFolder(path: string): Promise<void> {
        // Cannot open folders from web browser
        console.warn('OpenFolder not available in web mode');
    }

    async SelectFolder(defaultPath: string): Promise<string> {
        // Cannot select folders from web browser
        console.warn('SelectFolder not available in web mode');
        return defaultPath;
    }

    async SelectFile(): Promise<string> {
        // Cannot select files from web browser in the same way
        console.warn('SelectFile not available in web mode');
        return '';
    }

    // ==================== Analysis ====================

    async AnalyzeTrack(filePath: string): Promise<string> {
        const data = await this.fetch<any>('/api/analysis/track', {
            method: 'POST',
            body: JSON.stringify({ file_path: filePath }),
        });
        return JSON.stringify(data);
    }

    async AnalyzeMultipleTracks(filePaths: string[]): Promise<string> {
        // Would need server endpoint
        console.warn('AnalyzeMultipleTracks not implemented');
        return '[]';
    }

    // ==================== Lyrics & Covers ====================

    async DownloadLyrics(req: any): Promise<any> {
        // These could be added as API endpoints if needed
        console.warn('DownloadLyrics would need API endpoint');
        return { success: false };
    }

    async DownloadCover(req: any): Promise<any> {
        console.warn('DownloadCover would need API endpoint');
        return { success: false };
    }

    // ==================== Other ====================

    async CreateM3U8File(
        m3u8Name: string,
        outputDir: string,
        filePaths: string[]
    ): Promise<void> {
        // Would need server endpoint
        console.warn('CreateM3U8File not implemented');
    }

    async Quit(): Promise<void> {
        // Not applicable in web mode
        window.close();
    }
}

// Create singleton instance
export const apiClient = new ApiClient();

// Export for compatibility with Wails imports
export default apiClient;
