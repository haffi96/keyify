#!/bin/bash

STAGE=$1
if [ "$2" == "--local" ]; then
    LOCAL=true
else
    LOCAL=false
fi
PROJECT_ROOT=$(realpath $(dirname $0)/..)

function log() {
    local level="$1"
    shift
    printf "[${level}] $(date '+%Y-%m-%d %H:%M:%S') - $STAGE: %s\n" "$@" >&2
}

function preDeploy() {
    log INFO "Building and zipping Go Lambda functions..."
    mkdir -p dist

    LAMBDA_DIR=$(realpath src/lambdas)
    DIST_DIR=$(realpath dist)

    for dir in "$LAMBDA_DIR"/**/*.go; do
        dir_name=$(dirname "$dir")
        lambda="${dir_name##*/lambdas/}"

        cd "$dir_name"
        log INFO "Building Lambda function: $lambda"

        GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o $DIST_DIR/$lambda/bootstrap handler.go && cd $DIST_DIR/$lambda && zip function.zip bootstrap
    done
}

function cdkDeploy() {
    preDeploy

    log INFO "Deploying CDK stack..."
    cd $PROJECT_ROOT
    if [ "$LOCAL" = "true" ]; then
        log INFO "Deploying locally..."

        log INFO "Bootstrapping localstack..."
        CDK_ACCOUNT=000000000000 cdklocal bootstrap

        log INFO "Deploying CDK stack to localstack..."
        CDK_ACCOUNT=000000000000 cdklocal deploy dev-env/ApiKeyServiceStack --require-approval never
    else
        cdk deploy $STAGE-env/ApiKeyServiceStack --require-approval never
    fi

    log INFO "CDK stack deployed successfully!"

    log INFO "Cleaning up..."
    rm -rf dist
}

cdkDeploy
