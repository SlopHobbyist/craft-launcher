#!/bin/sh

# This script generates the manifest.json file on container startup
# It will be executed automatically by nginx's entrypoint system

MODPACK_DIR="/usr/share/nginx/html/files"
MANIFEST_FILE="/usr/share/nginx/html/manifest.json"
VERSION_FILE="/usr/share/nginx/html/files/.version"
OVERRIDES_FILE="/usr/share/nginx/html/files/.manifest_overrides"

# Create default overrides file if it doesn't exist
if [ ! -f "$OVERRIDES_FILE" ] && [ -d "$MODPACK_DIR" ]; then
    echo "Creating default .manifest_overrides..."
    cat > "$OVERRIDES_FILE" <<EOF
# Manifest Overrides Configuration
# Format: <glob_pattern> <override_boolean>
# Patterns are matched against the file path relative to the files directory.
# Rules are processed in order, last match wins.

# Default non-overridden user configuration files
options.txt false
optionsof.txt false
optionsshaders.txt false
servers.dat false
usercache.json false
EOF
fi

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
    
    # Skip if it's the version file or overrides file
    if [ "$REL_PATH" = ".version" ] || [ "$REL_PATH" = ".manifest_overrides" ]; then
        continue
    fi
    
    # Get file size
    SIZE=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null)
    
    # Calculate SHA256 checksum
    CHECKSUM=$(sha256sum "$file" | awk '{print $1}')
    
    # Determine override status (default true)
    OVERRIDE="true"
    
    # Check overrides file if it exists
    if [ -f "$OVERRIDES_FILE" ]; then
        # Read file line by line
        # We use a file descriptor other than 0 (stdin) to avoid conflict with the outer loop
        while IFS=' ' read -r pattern value <&3; do
            # Skip comments and empty lines
            case "$pattern" in 
                \#*|"") continue ;; 
            esac
            
            # Check if file matches pattern
            # We match against the relative path
            case "$REL_PATH" in
                $pattern) 
                    # Normalize boolean value
                    case "$value" in
                        [Tt][Rr][Uu][Ee]) OVERRIDE="true" ;;
                        [Ff][Aa][Ll][Ss][Ee]) OVERRIDE="false" ;;
                    esac
                    ;;
            esac
        done 3< "$OVERRIDES_FILE"
    else
        # Fallback to legacy hardcoded defaults for backward compatibility if file is missing (though we create it above)
        FILENAME=$(basename "$file")
        case "$FILENAME" in
            options.txt|optionsof.txt|optionsshaders.txt|servers.dat|usercache.json)
                OVERRIDE="false"
                ;;
        esac
    fi
    
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
