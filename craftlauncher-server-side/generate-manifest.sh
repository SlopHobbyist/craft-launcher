#!/bin/sh

# This script generates the manifest.json file on container startup
# It will be executed automatically by nginx's entrypoint system

MODPACK_DIR="/usr/share/nginx/html/files"
MANIFEST_FILE="/usr/share/nginx/html/manifest.json"
VERSION_FILE="/usr/share/nginx/html/files/.version"

# Clean up temp and system files
echo "Cleaning up temp and system files..."
find "$MODPACK_DIR" -type f \( \
    -name ".DS_Store" -o \
    -name "Thumbs.db" -o \
    -name "desktop.ini" -o \
    -name "*.tmp" -o \
    -name "*.temp" -o \
    -name "*.swp" -o \
    -name "*.swo" -o \
    -name "*~" -o \
    -name "._.DS_Store" -o \
    -name "._*" \
\) -delete
echo "Cleanup complete."

# Read version from .version file, default to 1 if not exists
if [ -f "$VERSION_FILE" ]; then
    VERSION=$(cat "$VERSION_FILE")
else
    VERSION=1
fi

echo "Generating manifest.json for version $VERSION..."

# Start JSON
echo "{" > "$MANIFEST_FILE"
echo "  \"version\": $VERSION," >> "$MANIFEST_FILE"
echo "  \"files\": [" >> "$MANIFEST_FILE"

# Generate file list with checksums
FIRST=true
find "$MODPACK_DIR" -type f ! -name ".version" | sort | while read -r file; do
    # Get relative path
    REL_PATH=$(echo "$file" | sed "s|$MODPACK_DIR/||")
    
    # Skip if it's the version file
    if [ "$REL_PATH" = ".version" ]; then
        continue
    fi
    
    # Get file size
    SIZE=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null)
    
    # Calculate SHA256 checksum
    CHECKSUM=$(sha256sum "$file" | awk '{print $1}')
    
    # Determine override status (default true)
    OVERRIDE="true"
    FILENAME=$(basename "$file")
    
    # Files that should NOT be overridden if they exist (user configs)
    case "$FILENAME" in
        options.txt|optionsof.txt|optionsshaders.txt|servers.dat|usercache.json)
            OVERRIDE="false"
            ;;
    esac
    
    # Add comma before each entry except the first
    if [ "$FIRST" = true ]; then
        FIRST=false
    else
        echo "," >> "$MANIFEST_FILE"
    fi
    
    # Write JSON entry (without trailing comma)
    printf "    {\n      \"path\": \"%s\",\n      \"size\": %s,\n      \"checksum\": \"%s\",\n      \"override\": %s\n    }" "$REL_PATH" "$SIZE" "$CHECKSUM" "$OVERRIDE" >> "$MANIFEST_FILE"
done

# Close JSON
echo "" >> "$MANIFEST_FILE"
echo "  ]" >> "$MANIFEST_FILE"
echo "}" >> "$MANIFEST_FILE"

echo "Manifest generated with $(grep -c '\"path\"' "$MANIFEST_FILE") files"
chmod 644 "$MANIFEST_FILE"
