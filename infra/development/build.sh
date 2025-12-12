#!/bin/bash

# Default values
SERVICE=""
MODE="${mode:-}"

echo "build script invoked..."

# Function to display usage
usage() {
    echo "Usage: $0 --service <service_name>"
    echo "Builds the specified service with environment-specific settings"
    exit 1
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --service)
            SERVICE="$2"
            shift 2
            ;;
        --help)
            usage
            ;;
        *)
            echo "Unknown option: $1"
            usage
            ;;
    esac
done

# Validate that service name is provided
if [[ -z "$SERVICE" ]]; then
    echo "Error: --service flag is required"
    usage
fi

# Check if mode is set to development
if [[ "$MODE" == "development" ]]; then
    echo "Development mode detected. Building for Linux..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "build/${SERVICE}" "./services/${SERVICE}/cmd/main.go"
    
    # Check if build was successful
    if [[ $? -eq 0 ]]; then
        echo "Build successful: build/${SERVICE}"
    else
        echo "Build failed!"
        exit 1
    fi
fi