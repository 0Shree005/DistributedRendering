#!/bin/bash

echo -n "Enter the save directory: "
read fileDir

fileDir=$(eval echo "$fileDir")

if [ ! -d "$fileDir" ]; then
    echo "Directory does not exist: $fileDir"
    exit 1
fi

nohup blender > /dev/null 2>&1 &
go run ~/code/DistributedRendering/client/client.go -fileDir="$fileDir"
