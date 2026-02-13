// Compatibility layer for Wails runtime functions
// Provides no-op or alternative implementations for web mode

/**
 * Quit function - closes window in web mode
 */
export const Quit = (): void => {
    window.close();
};

/**
 * EventsOn - replaced by WebSocket client
 * This is just for TypeScript compatibility
 */
export const EventsOn = (eventName: string, callback: (data: any) => void): void => {
    console.warn('EventsOn called - use WebSocket client instead');
};

/**
 * EventsOff - replaced by WebSocket client
 */
export const EventsOff = (eventName: string, callback?: (data: any) => void): void => {
    console.warn('EventsOff called - use WebSocket client instead');
};

/**
 * Window management functions - not applicable in web mode
 */
export const WindowMinimise = (): void => {
    console.warn('WindowMinimise not available in web mode');
};

export const WindowMaximise = (): void => {
    console.warn('WindowMaximise not available in web mode');
};

export const WindowToggleMaximise = (): void => {
    console.warn('WindowToggleMaximise not available in web mode');
};

export const WindowClose = (): void => {
    window.close();
};

/**
 * File drop functions - HTML5 drag & drop for web
 */
export const OnFileDrop = (callback: (x: number, y: number, paths: string[]) => void): void => {
    console.warn('OnFileDrop - not yet implemented for web mode');
    // TODO: Implement HTML5 drag & drop
};

export const OnFileDropOff = (): void => {
    console.warn('OnFileDropOff - not yet implemented for web mode');
    // TODO: Remove HTML5 drag & drop listeners
};

// Re-export everything as default for compatibility
export default {
    Quit,
    EventsOn,
    EventsOff,
    WindowMinimise,
    WindowMaximise,
    WindowToggleMaximise,
    WindowClose,
    OnFileDrop,
    OnFileDropOff,
};
