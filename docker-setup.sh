#!/usr/bin/env bash


VOLUME_NAME=my-vol
IMAGE_NAME=sdkapp

# Stop and remove any running containers
docker stop $(docker ps -a -q) 2>/dev/null || echo "No containers to stop"
docker rm $(docker ps -a -q) 2>/dev/null || echo "No containers to remove"

# Remove volume and recreate it
docker volume rm ${VOLUME_NAME} 2>/dev/null || echo "No volumes to remove"

# Create a new volume
docker volume create ${VOLUME_NAME}


# Build local docker container
docker build . -t ${IMAGE_NAME}

echo docker run -it --mount source=${VOLUME_NAME},target=/root ${IMAGE_NAME} /bin/bash
