#!/bin/bash

echo -n "Enter the save directory: "
read dirPath

dirPath=$(eval echo "$dirPath")

if [ ! -d "$dirPath" ]; then
    echo "Directory does not exist: $dirPath"
    exit 1
fi

nohup blender > /dev/null 2>&1 &
go run ~/code/DistributedRendering/main.go -dirPath="$dirPath"
