#!/usr/bin/env bash

set -eo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
cd "${DIR}"

function clean_up() {
  docker stop test-app >/dev/null 2>&1 || true
  docker network rm test-app-bridge >/dev/null 2>&1 || true
}

trap clean_up EXIT

docker network create -d bridge test-app-bridge

docker run -d --rm --name test-app --net=test-app-bridge vhndaree/test-app:develop

sleep 5 # This could fail some time if docker took more than 5 seconds to start

docker run --rm --name postman-newman -v "${DIR}/integration-test-demo.postman_collection.json:/postman_collection.json" --net=test-app-bridge postman/newman run --env-var "host=test-app" /postman_collection.json
