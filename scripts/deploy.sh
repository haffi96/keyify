#!/bin/bash

STAGE=$1 # Get stage from script input arguments

PROJECT_ROOT=$(realpath $(dirname $0)/..) # Get absolute path to root directory

function preDeploy() {
    echo "Building and zipping Go Lambda functions..."
    mkdir -p dist # Create assets directory if it doesn't exist

    LAMBDA_DIR=$(realpath apikeyservice/lambdas) # Get absolute path to root directory
    DIST_DIR=$(realpath dist)                    # Get absolute path to assets directory

    for dir in "$LAMBDA_DIR"/**/*.go; do
        # Extract directory name using dirname:
        dir_name=$(dirname "$dir")

        lambda="${dir_name##*/lambdas/}"

        # Change directory using absolute path:
        cd "$dir_name"

        # Build and zip using absolute paths:
        GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o $DIST_DIR/$lambda/bootstrap main.go && cd $DIST_DIR/$lambda && zip function.zip bootstrap
    done
}

function cdkDeploy() {
    preDeploy

    echo "Deploying CDK stack..."
    cd $PROJECT_ROOT
    cdk deploy $STAGE-env/ApiKeyServiceStack --require-approval never

    echo "CDK stack deployed successfully!"

    echo "Cleaning up..."
    rm -rf dist
}

cdkDeploy
