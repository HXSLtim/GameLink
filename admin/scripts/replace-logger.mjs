/**
 * Batch script to replace console.log/warn/error with logger utility
 * Run with: node scripts/replace-logger.mjs
 */

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Configuration
const SRC_DIR = path.join(__dirname, '../src');

// Track statistics
let stats = {
    filesProcessed: 0,
    filesModified: 0,
    consoleReplaced: {
        log: 0,
        warn: 0,
        error: 0,
        debug: 0,
        info: 0,
    },
    importsAdded: 0,
    filesSkipped: 0,
};

/**
 * Recursively find all TypeScript files
 */
function findFiles(dir, extensions = ['.ts', '.tsx']) {
    let files = [];
    const entries = fs.readdirSync(dir, { withFileTypes: true });

    for (const entry of entries) {
        const fullPath = path.join(dir, entry.name);

        if (entry.isDirectory()) {
            // Skip node_modules and hidden directories
            if (entry.name !== 'node_modules' && !entry.name.startsWith('.')) {
                files = files.concat(findFiles(fullPath, extensions));
            }
        } else if (entry.isFile()) {
            const ext = path.extname(entry.name);
            if (extensions.includes(ext)) {
                files.push(fullPath);
            }
        }
    }

    return files;
}

/**
 * Check if a file already has the logger import
 */
function hasLoggerImport(content) {
    return /import\s*{[^}]*logger[^}]*}\s*from\s*['"]@\/utils\/logger['"]/.test(content) ||
           /import\s*{[^}]*logger[^}]*}\s*from\s*['"]\.\/.*logger['"]/.test(content) ||
           /import\s*logger\s*from\s*['"]@\/utils\/logger['"]/.test(content);
}

/**
 * Add logger import to file content
 */
function addLoggerImport(content) {
    // Find the last import statement
    const importRegex = /^import\s+.*from\s+['"].*['"];?\s*$/gm;
    const imports = content.match(importRegex);

    if (!imports || imports.length === 0) {
        return content; // No imports found, skip
    }

    const lastImport = imports[imports.length - 1];
    const lastIndex = content.lastIndexOf(lastImport);
    const insertIndex = lastIndex + lastImport.length;

    // Insert logger import after the last import
    const before = content.substring(0, insertIndex);
    const after = content.substring(insertIndex);

    return `${before}\nimport { logger } from '@/utils/logger';${after}`;
}

/**
 * Replace console calls with logger calls
 */
function replaceConsoleCalls(content, filePath) {
    let modified = false;
    let replacements = 0;

    // Replace console.error → logger.error
    const errorMatches = content.match(/console\.error\(/g);
    if (errorMatches) {
        stats.consoleReplaced.error += errorMatches.length;
        replacements += errorMatches.length;
    }
    content = content.replace(/console\.error\(/g, 'logger.error(');
    if (errorMatches) modified = true;

    // Replace console.warn → logger.warn
    const warnMatches = content.match(/console\.warn\(/g);
    if (warnMatches) {
        stats.consoleReplaced.warn += warnMatches.length;
        replacements += warnMatches.length;
    }
    content = content.replace(/console\.warn\(/g, 'logger.warn(');
    if (warnMatches) modified = true;

    // Replace console.log → logger.info
    const logMatches = content.match(/console\.log\(/g);
    if (logMatches) {
        stats.consoleReplaced.log += logMatches.length;
        replacements += logMatches.length;
    }
    content = content.replace(/console\.log\(/g, 'logger.info(');
    if (logMatches) modified = true;

    // Replace console.debug → logger.debug
    const debugMatches = content.match(/console\.debug\(/g);
    if (debugMatches) {
        stats.consoleReplaced.debug += debugMatches.length;
        replacements += debugMatches.length;
    }
    content = content.replace(/console\.debug\(/g, 'logger.debug(');
    if (debugMatches) modified = true;

    // Replace console.info → logger.info
    const infoMatches = content.match(/console\.info\(/g);
    if (infoMatches) {
        stats.consoleReplaced.info += infoMatches.length;
        replacements += infoMatches.length;
    }
    content = content.replace(/console\.info\(/g, 'logger.info(');
    if (infoMatches) modified = true;

    // Add logger import if needed and file was modified
    if (modified && !hasLoggerImport(content)) {
        content = addLoggerImport(content);
        stats.importsAdded++;
    }

    return { content, modified, replacements };
}

/**
 * Process a single file
 */
function processFile(filePath) {
    stats.filesProcessed++;

    try {
        let content = fs.readFileSync(filePath, 'utf-8');

        // Skip if no console calls
        if (!content.includes('console.')) {
            stats.filesSkipped++;
            return;
        }

        // Skip documentation files
        if (filePath.endsWith('.md')) {
            stats.filesSkipped++;
            return;
        }

        // Replace console calls
        const { content: newContent, modified, replacements } = replaceConsoleCalls(content, filePath);

        if (modified) {
            fs.writeFileSync(filePath, newContent, 'utf-8');
            stats.filesModified++;
            console.log(`✓ Modified: ${path.relative(SRC_DIR, filePath)} (${replacements} replacements)`);
        } else {
            stats.filesSkipped++;
        }
    } catch (error) {
        console.error(`✗ Error processing ${filePath}:`, error.message);
    }
}

/**
 * Main function
 */
async function main() {
    console.log('🔍 Searching for TypeScript files with console calls...');
    console.log(`📁 Source directory: ${SRC_DIR}\n`);

    // Find all TypeScript files
    const files = findFiles(SRC_DIR, ['.ts', '.tsx']);

    console.log(`📊 Found ${files.length} total files\n`);

    // Process each file
    for (const file of files) {
        processFile(file);
    }

    // Print statistics
    console.log('\n' + '='.repeat(60));
    console.log('📈 Replacement Statistics:');
    console.log('='.repeat(60));
    console.log(`Files processed:    ${stats.filesProcessed}`);
    console.log(`Files modified:     ${stats.filesModified}`);
    console.log(`Files skipped:      ${stats.filesSkipped}`);
    console.log(`Imports added:      ${stats.importsAdded}`);
    console.log('\nReplacements by type:');
    console.log(`  console.log   → logger.info:  ${stats.consoleReplaced.log}`);
    console.log(`  console.warn  → logger.warn:  ${stats.consoleReplaced.warn}`);
    console.log(`  console.error → logger.error: ${stats.consoleReplaced.error}`);
    console.log(`  console.debug → logger.debug: ${stats.consoleReplaced.debug}`);
    console.log(`  console.info  → logger.info:  ${stats.consoleReplaced.info}`);
    console.log(`  Total:                        ${Object.values(stats.consoleReplaced).reduce((a, b) => a + b, 0)}`);
    console.log('='.repeat(60));

    if (stats.filesModified > 0) {
        console.log('\n✅ Done! Please review the changes and run tests.');
    } else {
        console.log('\nℹ️  No files needed modification.');
    }
}

// Run the script
main().catch(error => {
    console.error('Fatal error:', error);
    process.exit(1);
});
