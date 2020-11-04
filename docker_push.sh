#!/bin/bash
# echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker login -u "$DOCKER_USERNAME" -p "$DOCKER_PASSWORD"

make docker-push IMG=${REPO}:${COMMIT}
