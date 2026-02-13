/**
 * Backend API Types (replacing Wails models)
 * These types match the API responses from the Go backend
 */

// Search response types
export interface SearchResultItem {
    id: string;
    name: string;
    external_urls: string;
    images?: string;
    artists?: string;
    duration_ms?: number;
    is_explicit?: boolean;
    release_date?: string;
    owner?: string;
}

export interface SearchResponse {
    tracks: SearchResultItem[];
    albums: SearchResultItem[];
    artists: SearchResultItem[];
    playlists: SearchResultItem[];
}

// Download queue types
export interface DownloadItem {
    id: string;
    title: string;
    artists: string;
    album?: string;
    status: "pending" | "downloading" | "completed" | "failed";
    progress: number;
    error?: string;
    url: string;
    quality: string;
    downloader: string;
}

export interface DownloadQueueResponse {
    items: DownloadItem[];
    total: number;
    pending: number;
    downloading: number;
    completed: number;
    failed: number;
}
