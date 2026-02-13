#!/usr/bin/env node

/**
 * Migration Script: Wails to HTTP API
 * This script automatically migrates all Wails imports to the new HTTP API client
 */

const fs = require('fs');
const path = require('path');

const FRONTEND_SRC = path.join(__dirname, 'frontend', 'src');

// Migration rules
const MIGRATIONS = [
    // Remove Wails App imports
    {
        pattern: /import\s+\{[^}]+\}\s+from\s+["']\.\.\/\.\.\/wailsjs\/go\/main\/App["'];?\s*\n/g,
        replacement: ''
    },
    {
        pattern: /import\s+\{[^}]+\}\s+from\s+["']@\/\.\.\/wailsjs\/go\/main\/App["'];?\s*\n/g,
        replacement: ''
    },

    // Remove Wails models imports
    {
        pattern: /import\s+\{\s*backend\s*\}\s+from\s+["']\.\.\/\.\.\/wailsjs\/go\/models["'];?\s*\n/g,
        replacement: ''
    },

    // Remove Wails runtime imports (but keep if it's runtime.ts - our compat layer)
    {
        pattern: /import\s+\{[^}]+\}\s+from\s+["']\.\.\/\.\.\/wailsjs\/runtime\/runtime["'];?\s*\n/g,
        replacement: (match, offset, string, filePath) => {
            // Don't replace if we're in runtime.ts itself
            if (filePath && filePath.includes('api/runtime.ts')) {
                return match;
            }
            return '';
        }
    },

    // Replace dynamic imports
    {
        pattern: /const\s+\{([^}]+)\}\s+=\s+await\s+import\(["']\.\.\/\.\.\/wailsjs\/go\/main\/App["']\);/g,
        replacement: 'const { $1 } = apiClient;'
    },

    // Replace backend.SearchResponse with SearchResponse
    {
        pattern: /backend\.SearchResponse/g,
        replacement: 'SearchResponse'
    },

    // Replace backend.DownloadItem with DownloadItem
    {
        pattern: /backend\.DownloadItem/g,
        replacement: 'DownloadItem'
    }
];

// Files that need specific imports added
const IMPORT_ADDITIONS = {
    needsApiClient: /SearchSpotify|GetDownloadQueue|GetDownloadProgress|AnalyzeTrack|CheckTrackAvailability|GetPreviewURL|SelectFolder|GetOSInfo|ClearCompletedDownloads|ClearAllDownloads|ExportFailedDownloads|AddToDownloadQueue|GetStreamingURLs|MarkDownloadItemFailed|CancelAllQueuedItems|ConvertAudio|SelectAudioFiles|SelectFile|UploadImageBytes|UploadImage|SelectImageVideo|GetDownloadHistory|ClearDownloadHistory|DeleteDownloadHistoryItem|GetFetchHistory|DeleteFetchHistoryItem|ClearFetchHistoryByType/,
    needsBackendTypes: /SearchResponse|DownloadItem|DownloadQueueResponse/,
    needsRuntime: /WindowMinimise|WindowToggleMaximise|Quit|OnFileDrop|OnFileDropOff/
};

function findTsFiles(dir, files = []) {
    const items = fs.readdirSync(dir);

    for (const item of items) {
        const fullPath = path.join(dir, item);
        const stat = fs.statSync(fullPath);

        if (stat.isDirectory()) {
            // Skip node_modules and dist
            if (item !== 'node_modules' && item !== 'dist' && item !== 'build') {
                findTsFiles(fullPath, files);
            }
        } else if (item.endsWith('.ts') || item.endsWith('.tsx')) {
            files.push(fullPath);
        }
    }

    return files;
}

function migrateFile(filePath) {
    let content = fs.readFileSync(filePath, 'utf8');
    const originalContent = content;
    let modified = false;

    // Apply all migration rules
    for (const migration of MIGRATIONS) {
        const beforeLength = content.length;
        if (typeof migration.replacement === 'function') {
            content = content.replace(migration.pattern, (match, ...args) => {
                return migration.replacement(match, ...args, filePath);
            });
        } else {
            content = content.replace(migration.pattern, migration.replacement);
        }
        if (content.length !== beforeLength) {
            modified = true;
        }
    }

    // Check if we need to add imports
    const needsApiClient = IMPORT_ADDITIONS.needsApiClient.test(content);
    const needsBackendTypes = IMPORT_ADDITIONS.needsBackendTypes.test(content);
    const needsRuntime = IMPORT_ADDITIONS.needsRuntime.test(content);

    // Only add imports if they don't already exist
    const hasApiClientImport = /import\s+\{[^}]*apiClient[^}]*\}\s+from/.test(content);
    const hasBackendTypesImport = /import\s+.*from\s+["']@\/types\/backend["']/.test(content);
    const hasRuntimeImport = /import\s+\{[^}]+\}\s+from\s+["']\.\.\/api\/runtime["']/.test(content) ||
        /import\s+\{[^}]+\}\s+from\s+["']\.\.\/\.\.\/api\/runtime["']/.test(content);

    // Find the last import statement
    const importRegex = /^import\s+.*?;?\s*$/gm;
    const imports = content.match(importRegex);

    if (imports && imports.length > 0) {
        const lastImport = imports[imports.length - 1];
        const lastImportIndex = content.lastIndexOf(lastImport) + lastImport.length;

        let newImports = [];

        if (needsApiClient && !hasApiClientImport) {
            // Determine correct path based on file location
            const relativePath = path.relative(path.dirname(filePath), path.join(FRONTEND_SRC, 'api', 'client.ts'));
            const importPath = relativePath.startsWith('.') ? relativePath.replace(/\.ts$/, '') : `./${relativePath.replace(/\.ts$/, '')}`;
            newImports.push(`import { apiClient } from "${importPath.replace(/\\/g, '/')}";`);
            modified = true;
        }

        if (needsBackendTypes && !hasBackendTypesImport) {
            newImports.push(`import type { SearchResponse, DownloadItem, DownloadQueueResponse } from "@/types/backend";`);
            modified = true;
        }

        if (needsRuntime && !hasRuntimeImport) {
            const relativePath = path.relative(path.dirname(filePath), path.join(FRONTEND_SRC, 'api', 'runtime.ts'));
            const importPath = relativePath.startsWith('.') ? relativePath.replace(/\.ts$/, '') : `./${relativePath.replace(/\.ts$/, '')}`;
            newImports.push(`import { WindowMinimise, WindowToggleMaximise, Quit, OnFileDrop, OnFileDropOff } from "${importPath.replace(/\\/g, '/')}";`);
            modified = true;
        }

        if (newImports.length > 0) {
            content = content.slice(0, lastImportIndex) + '\n' + newImports.join('\n') + content.slice(lastImportIndex);
        }
    }

    if (modified && content !== originalContent) {
        fs.writeFileSync(filePath, content, 'utf8');
        return true;
    }

    return false;
}

function main() {
    console.log('🔍 Finding TypeScript files...');
    const files = findTsFiles(FRONTEND_SRC);
    console.log(`📁 Found ${files.length} TypeScript files\n`);

    let migratedCount = 0;
    const migratedFiles = [];

    for (const file of files) {
        // Skip our new API files
        if (file.includes('/api/client.ts') ||
            file.includes('/api/websocket.ts') ||
            file.includes('/api/runtime.ts') ||
            file.includes('/types/backend.ts')) {
            continue;
        }

        const relativePath = path.relative(FRONTEND_SRC, file);

        try {
            if (migrateFile(file)) {
                console.log(`✅ Migrated: ${relativePath}`);
                migratedCount++;
                migratedFiles.push(relativePath);
            }
        } catch (error) {
            console.error(`❌ Error migrating ${relativePath}:`, error.message);
        }
    }

    console.log(`\n✨ Migration complete!`);
    console.log(`📊 Migrated ${migratedCount} files`);

    if (migratedFiles.length > 0) {
        console.log(`\n📝 Modified files:`);
        migratedFiles.forEach(f => console.log(`   - ${f}`));
    }
}

main();
