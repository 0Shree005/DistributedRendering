#!/bin/bash

fileName="$1"

fullFilePath="/home/biggiecheese/code/DistributedRendering/server/uploads/$fileName"

outputDir="/home/biggiecheese/code/DistributedRendering/server/renders"

outputFileName="$(basename "$fileName" .blend)_render"

if [[ ! -f "$fullFilePath" ]]; then
  echo "File not found: $fullFilePath"
  exit 1
fi

echo "Starting render for $fullFilePath..."
blender -b "$fullFilePath" -o "$outputDir/$outputFileName" -F PNG -x 1 -f 1
status=$?

if [[ $status -eq 0 ]]; then
  echo "Rendering completed for $fileName."
else
  echo "Rendering failed for $fileName (exit code $status)."
  exit $status
fi


# Open the rendered image with default image viewer
echo "Opening rendered image..."
xdg-open "${outputDir}/${outputFileName}0001.png"
