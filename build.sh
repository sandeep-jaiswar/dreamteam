#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status
set -u  # Treat unset variables as an error

# Variables
APP_NAME="dreamteam"
BUILD_DIR="./build"
BINARY_NAME="app"
DOCKER_IMAGE="dreamteam:latest"
LINT_BIN="./bin/golangci-lint"

# Colors for output
GREEN="\033[0;32m"
RED="\033[0;31m"
NC="\033[0m" # No Color

echo -e "${GREEN}Starting build script...${NC}"

# 1. Clean previous builds
echo -e "${GREEN}Cleaning up old build files...${NC}"
rm -rf ${BUILD_DIR}
mkdir -p ${BUILD_DIR}

# 2. Download and verify dependencies
echo -e "${GREEN}Downloading dependencies...${NC}"
go mod tidy
go mod verify

# 3. Code formatting
echo -e "${GREEN}Formatting code...${NC}"
go fmt ./...

# 4. Linting code
echo -e "${GREEN}Checking for golangci-lint...${NC}"
if ! [ -x "$(command -v golangci-lint)" ] && ! [ -x "${LINT_BIN}" ]; then
  echo -e "${RED}golangci-lint not found. Installing...${NC}"
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
  mkdir -p ./bin
  mv ./golangci-lint ./bin/
fi

echo -e "${GREEN}Running golangci-lint...${NC}"
if [ -f "${LINT_BIN}" ]; then
  ${LINT_BIN} run
else
  golangci-lint run
fi

# 5. Static analysis and tests
echo -e "${GREEN}Running static analysis and tests...${NC}"
go test ./... -cover -v

# 6. Build the binary
echo -e "${GREEN}Building the Go binary...${NC}"
GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/${BINARY_NAME} ./cmd/api/main.go

# 7. Copy configuration files
echo -e "${GREEN}Copying configuration files...${NC}"
cp -r ./configs ${BUILD_DIR}/

# 8. Optional: Build Docker image
if [ "${1:-}" == "docker" ]; then
  echo -e "${GREEN}Building Docker image...${NC}"
  docker build -t ${DOCKER_IMAGE} .
fi

echo -e "${GREEN}Build process completed successfully!${NC}"
