#!/bin/bash

set -e

ECR_REGISTRY=123456789.dkr.ecr.us-east-1.amazonaws.com
ECR_REPOSITORY=penguin-server
IMAGE_TAG="latest"
CONTAINER_NAME="penguin-server"
REGION=us-east-1

# Get the latest digest from ECR
LATEST_DIGEST=$(aws ecr describe-images \
  --repository-name $ECR_REPOSITORY \
  --region $REGION \
  --query "sort_by(imageDetails,& imagePushedAt)[-1].imageDigest" \
  --output text | tr -d '\r\n')

# Compose full image reference
FULL_IMAGE="$ECR_REGISTRY/$ECR_REPOSITORY@$LATEST_DIGEST"

# Get the image ID used by the container (if it exists)
IMAGE_ID=$(docker inspect --format='{{.Image}}' $CONTAINER_NAME 2>/dev/null || echo "")

# Resolve the digest for that image ID
CURRENT_DIGEST=""
if [ -n "$IMAGE_ID" ]; then
  CURRENT_DIGEST=$(docker inspect --format='{{range .RepoDigests}}{{println .}}{{end}}' "$IMAGE_ID" 2>/dev/null \
    | grep "$ECR_REPOSITORY@" | awk -F@ '{print $2}' | head -n 1 | tr -d '\r\n')
fi

echo "Latest digest:   $LATEST_DIGEST"
echo "Current digest:  $CURRENT_DIGEST"

if [ "$CURRENT_DIGEST" == "$LATEST_DIGEST" ]; then
  echo "The container is already running the latest image. No action needed."
  exit 0
fi

# Login to ECR
aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $ECR_REGISTRY

# Pull the image by digest
docker pull "$FULL_IMAGE"

echo "Updating container to the latest image..."
docker stop $CONTAINER_NAME 2>/dev/null || true
docker rm $CONTAINER_NAME 2>/dev/null || true

docker run -d --name $CONTAINER_NAME -p 7777:7777 "$FULL_IMAGE"
