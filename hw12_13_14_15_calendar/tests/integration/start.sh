#!/bin/bash

docker-compose --env-file ./deployments/.env.tests -f ./deployments/docker-compose.yaml run --build --rm tests
exit_code=$?

docker-compose -f ./deployments/docker-compose.yaml down

echo "Exit code: $exit_code"
exit $exit_code
