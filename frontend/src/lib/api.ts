import type { SpotifyMetadataResponse, DownloadRequest, DownloadResponse, HealthResponse, LyricsDownloadRequest, LyricsDownloadResponse, CoverDownloadRequest, CoverDownloadResponse, HeaderDownloadRequest, HeaderDownloadResponse, GalleryImageDownloadRequest, GalleryImageDownloadResponse, AvatarDownloadRequest, AvatarDownloadResponse, } from "@/types/api";
import { apiClient } from "../api/client";

// Re-export API client methods for compatibility
export const GetSpotifyMetadata = (req: any) => apiClient.GetSpotifyMetadata(req);
export const DownloadTrack = (req: any) => apiClient.DownloadTrack(req);
export const DownloadLyrics = (req: any) => apiClient.DownloadLyrics(req);
export const DownloadCover = (req: any) => apiClient.DownloadCover(req);
// Note: Image download functions would need server endpoints
export const DownloadHeader = (req: any) => Promise.resolve({ success: false });
export const DownloadGalleryImage = (req: any) => Promise.resolve({ success: false });
export const DownloadAvatar = (req: any) => Promise.resolve({ success: false });
export async function fetchSpotifyMetadata(url: string, batch: boolean = true, delay: number = 1.0, timeout: number = 300.0): Promise<SpotifyMetadataResponse> {
    const req = new main.SpotifyMetadataRequest({
        url,
        batch,
        delay,
        timeout,
    });
    const jsonString = await GetSpotifyMetadata(req);
    return JSON.parse(jsonString);
}
export async function downloadTrack(request: DownloadRequest): Promise<DownloadResponse> {
    const req = new main.DownloadRequest(request);
    return await DownloadTrack(req);
}
export async function checkHealth(): Promise<HealthResponse> {
    return {
        status: "ok",
        time: new Date().toISOString(),
    };
}
export async function downloadLyrics(request: LyricsDownloadRequest): Promise<LyricsDownloadResponse> {
    const req = new main.LyricsDownloadRequest(request);
    return await DownloadLyrics(req);
}
export async function downloadCover(request: CoverDownloadRequest): Promise<CoverDownloadResponse> {
    const req = new main.CoverDownloadRequest(request);
    return await DownloadCover(req);
}
export async function downloadHeader(request: HeaderDownloadRequest): Promise<HeaderDownloadResponse> {
    const req = new main.HeaderDownloadRequest(request);
    return await DownloadHeader(req);
}
export async function downloadGalleryImage(request: GalleryImageDownloadRequest): Promise<GalleryImageDownloadResponse> {
    const req = new main.GalleryImageDownloadRequest(request);
    return await DownloadGalleryImage(req);
}
export async function downloadAvatar(request: AvatarDownloadRequest): Promise<AvatarDownloadResponse> {
    const req = new main.AvatarDownloadRequest(request);
    return await DownloadAvatar(req);
}
