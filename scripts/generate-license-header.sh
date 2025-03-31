#!/bin/bash

# Set the current year for the copyright notice
YEAR=$(date +"%Y")
# Set your name as it should appear in the copyright notice
AUTHOR="Praveen Kumar"

# Define file types to process
FILE_TYPES=("tsx" "ts" "go")

# Create a temporary file for the license header
TMP_HEADER=$(mktemp)

# Write the license header to the temporary file
cat > "$TMP_HEADER" << EOF
/*
 * Copyright $YEAR $AUTHOR
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

EOF

# Function to check if file already has a license header
has_license() {
    grep -q "Licensed under the Apache License" "$1"
    return $?
}

# Counter for files processed
FILES_PROCESSED=0
FILES_UPDATED=0

# Build the file extension regex pattern from the FILE_TYPES array
FILE_PATTERN=""
FIND_PATTERN=""
for i in "${!FILE_TYPES[@]}"; do
    ext="${FILE_TYPES[$i]}"
    # For grep pattern
    if [ $i -eq 0 ]; then
        FILE_PATTERN="\\.$ext\$"
    else
        FILE_PATTERN="$FILE_PATTERN|\\.$ext\$"
    fi
    
    # For find pattern
    if [ $i -eq 0 ]; then
        FIND_PATTERN="-name \"*.$ext\""
    else
        FIND_PATTERN="$FIND_PATTERN -o -name \"*.$ext\""
    fi
done

# Use git ls-files to only find files that are not in .gitignore
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    # We're in a git repository, use git ls-files to get tracked files and untracked files not in .gitignore
    echo "Using git to find files (respecting .gitignore)..."
    
    # Get list of all non-ignored files (both tracked and untracked)
    FILE_LIST=$(git ls-files --cached --others --exclude-standard | grep -E "$FILE_PATTERN")
else
    # If not in a git repo, use find but manually filter using .gitignore if it exists
    echo "Not a git repository, using find command..."
    
    # Start with empty exclude pattern
    EXCLUDE_PATTERN=""
    
    # If .gitignore exists, build exclude pattern from it
    if [ -f ".gitignore" ]; then
        echo "Using patterns from .gitignore..."
        
        # Common patterns to always exclude (including node_modules)
        EXCLUDE_PATTERN="-path './node_modules*' -prune -o"
        
        # Add patterns from .gitignore (handling basic patterns only)
        while IFS= read -r pattern; do
            # Skip empty lines and comments
            if [[ -z "$pattern" || "$pattern" == \#* ]]; then
                continue
            fi
            
            # Handle basic directory patterns (those ending with /)
            if [[ "$pattern" == */ ]]; then
                pattern="${pattern%/}"
                EXCLUDE_PATTERN="$EXCLUDE_PATTERN -path './$pattern*' -prune -o"
            # Handle basic file and wildcard patterns
            else
                EXCLUDE_PATTERN="$EXCLUDE_PATTERN -path './$pattern' -prune -o"
            fi
        done < .gitignore
    else
        echo "No .gitignore found, excluding node_modules by default..."
        EXCLUDE_PATTERN="-path './node_modules*' -prune -o"
    fi
    
    # Find files with the exclude pattern - use eval to properly handle the constructed pattern
    eval "FILE_LIST=\$(find . \$EXCLUDE_PATTERN -type f \( $FIND_PATTERN \) -print)"
fi

# Process each file
echo "$FILE_LIST" | while read -r file; do
    # Skip empty lines
    [ -z "$file" ] && continue
    
    FILES_PROCESSED=$((FILES_PROCESSED + 1))
    
    # Skip files that already have a license header
    if has_license "$file"; then
        echo "Skipping file (already has license): $file"
        continue
    fi

    echo "Adding license header to: $file"
    
    # Create a temporary file with the license header followed by the original content
    TMP_FILE=$(mktemp)
    cat "$TMP_HEADER" "$file" > "$TMP_FILE"
    
    # Replace the original file with the new content
    mv "$TMP_FILE" "$file"
    
    FILES_UPDATED=$((FILES_UPDATED + 1))
done

# Clean up the temporary header file
rm "$TMP_HEADER"

echo "Summary:"
echo "Files processed: $FILES_PROCESSED"
echo "Files updated: $FILES_UPDATED"
echo "Done!"