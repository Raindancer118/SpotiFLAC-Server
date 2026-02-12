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

// Re-export everything as default for compatibility
export default {
    Quit,
    EventsOn,
    EventsOff,
    WindowMinimise,
    WindowMaximise,
    WindowToggleMaximise,
    WindowClose,
};
