#!/bin/bash

# Copyright (c) 2025 A Bit of Help, Inc.

# Generate SVG files from PlantUML source files
echo "Generating SVG files from PlantUML source files..."
plantuml -tsvg DOCS/diagrams/source/*.puml

# Move SVG files to the svg directory
echo "Moving SVG files to the svg directory..."
mv DOCS/diagrams/source/*.svg DOCS/diagrams/svg/

echo "SVG files generated and moved successfully."