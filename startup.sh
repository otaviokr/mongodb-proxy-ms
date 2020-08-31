#!/bin/bash

GOPATH="/home/otaviokr/go:/go"
USERHOME=${HOME}
CURRENT_DIR="${PWD}"
SWAGGER_OUTPUT="swagger.json"

# Touching the YAML file to avoid having root as author!
touch "${SWAGGER_OUTPUT}"

echo "Generating the YAML file for swagger..."
docker run --rm -it \
    -e GOPATH="${GOPATH}" \
    -v ${HOME}:${HOME} \
    -w ${CURRENT_DIR} \
    quay.io/goswagger/swagger \
    generate spec -o ${SWAGGER_OUTPUT} --scan-models

echo "Starting docker containers..."
docker-compose up