#!/bin/bash
set -e

echo "Building docker image"
docker build -t terraforming-mars:latest .

echo "Stopping (posibly) running container"
docker stop tm-server 2>/dev/null || true
docker rm tm-server 2>/dev/null || true

if [ ! -d ./data ];
then
  echo "Creating data directory"
  mkdir "./data"
  sudo chown 65534:65534 "./data"
  sudo chmod 777 "./data"
fi

echo "Starting container"
docker run -d \
  --name tm-server \
  --restart unless-stopped \
  -p 127.0.0.1:8080:8080 \
  -v ./data:/data \
  --read-only \
  --tmpfs /tmp:rw,noexec,nosuid,size=64m \
  --memory=256m \
  terraforming-mars:latest
