#!/bin/bash

ECR_REGISTRY=123456789012.dkr.ecr.us-east-1.amazonaws.com
ECR_REPOSITORY=penguin-server
IMAGE_TAG="latest"
CONTAINER_NAME="penguin-server"
REGION=us-east-1

aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $ECR_REGISTRY

# Get the latest image digest from ECR
LATEST_DIGEST=$(aws ecr describe-images --repository-name $ECR_REPOSITORY --region $REGION --query "imageDetails[?imageTags[?contains(@, '$IMAGE_TAG')]].imageDigest" --output text)

# Get the currently running container's image digest
CURRENT_DIGEST=$(docker inspect --format='{{index .RepoDigests 0}}' $CONTAINER_NAME 2>/dev/null | awk -F'@' '{print $2}')

# Compare digests
if [ "$LATEST_DIGEST" == "$CURRENT_DIGEST" ]; then
  echo "The container is already running the latest image. No action needed."
  exit 0
fi

echo "New image detected. Updating the container..."

docker pull $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG

docker stop $CONTAINER_NAME 
docker rm $CONTAINER_NAME 

docker run -d --name $CONTAINER_NAME -p 7777:7777 $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG
