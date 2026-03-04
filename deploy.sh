#!/bin/bash
set -euo pipefail

echo "Building arm image..."
docker buildx build --platform linux/arm/v7 -t go-alexa-api:latest -o type=docker,dest=go-alexa-api.tar .

echo "Transferring to media server..."
scp go-alexa-api.tar docker-compose.yml .env pi@media-server.home:~/go-alexa-api/

echo "Deploying..."
ssh pi@media-server.home "cd ~/go-alexa-api && docker load -i go-alexa-api.tar && docker compose up -d"

echo "Done!"
