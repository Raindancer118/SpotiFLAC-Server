import type { SpotifyMetadataResponse, DownloadRequest, DownloadResponse, HealthResponse, LyricsDownloadRequest, LyricsDownloadResponse, CoverDownloadRequest, CoverDownloadResponse, HeaderDownloadRequest, HeaderDownloadResponse, GalleryImageDownloadRequest, GalleryImageDownloadResponse, AvatarDownloadRequest, AvatarDownloadResponse, } from "@/types/api";
import { apiClient } from "../api/client";

/**
 * Fetches Spotify metadata for a given URL
 * @param url - Spotify URL (track, album, playlist, or artist)
 * @param batch - Whether to batch fetch related items
 * @param delay - Delay between requests in seconds
 * @param timeout - Request timeout in seconds
 */
export async function fetchSpotifyMetadata(url: string, batch: boolean = true, delay: number = 1.0, timeout: number = 300.0): Promise<SpotifyMetadataResponse> {
    const req = {
        url,
        batch,
        delay,
        timeout,
    };
    const jsonString = await apiClient.GetSpotifyMetadata(req);
    return JSON.parse(jsonString);
}

/**
 * Initiates a track download
 */
export async function downloadTrack(request: DownloadRequest): Promise<DownloadResponse> {
    return await apiClient.DownloadTrack(request);
}

/**
 * Checks API health status
 */
export async function checkHealth(): Promise<HealthResponse> {
    return {
        status: "ok",
        time: new Date().toISOString(),
    };
}

/**
 * Downloads lyrics for a track
 */
export async function downloadLyrics(request: LyricsDownloadRequest): Promise<LyricsDownloadResponse> {
    return await apiClient.DownloadLyrics(request);
}

/**
 * Downloads album cover art
 */
export async function downloadCover(request: CoverDownloadRequest): Promise<CoverDownloadResponse> {
    return await apiClient.DownloadCover(request);
}

// Note: The following image download functions are not yet implemented on the server
export async function downloadHeader(request: HeaderDownloadRequest): Promise<HeaderDownloadResponse> {
    return { success: false, message: "Not implemented" };
}

export async function downloadGalleryImage(request: GalleryImageDownloadRequest): Promise<GalleryImageDownloadResponse> {
    return { success: false, message: "Not implemented" };
}

export async function downloadAvatar(request: AvatarDownloadRequest): Promise<AvatarDownloadResponse> {
    return { success: false, message: "Not implemented" };
}
