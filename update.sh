# Updates the image with the current changes and pushes it to the registry.
# Usage: ./update.sh

#!/bin/bash

# Check if the image tag argument is provided
if [ -z "$1" ]; then
    echo "You need to provide an image tag argument."
    exit 1
fi

# Set the tag name and Docker username based on the command line argument
TAG="$1"
DOCKER_USERNAME="reagancn"
IMAGE_NAME="telegram-gpt"

# Build the Docker image
echo "Setting up Docker buildx with custom builder..."
docker buildx use mybuilder

echo "Building Docker image..."
# docker buildx build --platform linux/amd64,linux/arm64 -t reagancn/telegram-gpt:0.0.8 --push .
docker buildx build --platform linux/amd64,linux/arm64 -t $DOCKER_USERNAME/$IMAGE_NAME:$TAG --push .

# Login to Docker Hub
echo "Logging into Docker Hub..."
docker login
