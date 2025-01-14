#!/bin/bash

# Get the file name from the first argument
fileName="$1"

# Full file path
fullFilePath="/home/biggiecheese/code/DistributedRendering/server/uploads/$fileName"

# Ensure the file exists
if [[ ! -f "$fullFilePath" ]]; then
  echo "File not found: $fullFilePath"
  exit 1
fi

# Render output file name
outputFileName="$(basename "$fileName" .blend)render"

# Run Blender in background mode and wait for it to finish
echo "Starting render for $fullFilePath..."
blender -b "$fullFilePath" -o "./renders/$outputFileName#" -F PNG -x 1 -f 1

# Check if Blender rendered successfully
if [[ $? -eq 0 ]]; then
  echo "Rendering completed for $fileName."
else
  echo "Rendering failed for $fileName."
  exit 1
fi



# Open the rendered image with feh
echo "Opening rendered image..."
feh --fullscreen "./renders/${outputFileName}0001.png"
