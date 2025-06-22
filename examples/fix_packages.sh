#!/bin/bash

# Function to rename package in a file
rename_package() {
    local file=$1
    local package_name=$2
    sed -i '' "s/^package main$/package ${package_name}/" "$file"
}

# Process each subdirectory in examples
for dir in /Users/mike/Go/src/github.com/abitofhelp/servicelib/examples/*/; do
    if [ -d "$dir" ]; then
        # Get directory name
        dirname=$(basename "$dir")
        # Process each Go file in the directory
        for file in "$dir"*.go; do
            if [ -f "$file" ]; then
                rename_package "$file" "example_${dirname}"
            fi
        done
    fi
done

# Process files in examples root directory
for file in /Users/mike/Go/src/github.com/abitofhelp/servicelib/examples/*.go; do
    if [ -f "$file" ]; then
        filename=$(basename "$file" .go)
        rename_package "$file" "example_${filename}"
    fi
done
