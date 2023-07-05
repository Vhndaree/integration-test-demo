#!/usr/bin/env bash

set -eo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
cd "${SCRIPT_DIR}"

IMAGE_TAG=$1

if [ -z "${IMAGE_TAG}" ]; then
  echo "No IMAGE_TAG set. Using default [develop]."
  export IMAGE_TAG="develop"
fi

function wait_for_container() {
  CONTAINER=$1
  MSG=$2
  i=0

  while ! docker logs "${CONTAINER}" 2>&1 | grep -m 1 "${MSG}"; do
    if [[ "$i" -gt 20 ]]; then
      echo "timeout"
      docker logs "${CONTAINER}" --tail=20
      exit 1
    fi
    sleep 5
    echo "waiting for [${MSG}] from ${CONTAINER}"
    i=$((i+1))
  done
  echo "get [${MSG}] from ${CONTAINER}"
}

function clean_up() {
  docker stop test-app || true
  docker network rm test-app-bridge || true
}

function ramp_up() {
  trap "clean_up" EXIT

  docker network create -d bridge test-app-bridge
  docker build -t test-app ../.
  docker run --rm --name test-app  --net=test-app-bridge -d vhndaree/test-app:develop
  wait_for_container "test-app" "Server starting on port 8848"
}

ramp_up

docker run --rm --name postman-newman -v "${SCRIPT_DIR}/integration-test-demo.postman_collection.json:/postman_collection.json" --net=test-app-bridge postman/newman run --env-var "host=test-app" /postman_collection.json
