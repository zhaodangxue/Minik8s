#!/bin/bash

build_dir=./build

deploy_dir=/opt/minik8s

execute_command() {
    local command=$1
    local description=$2

    echo "Executing: $description"
    $command
    if [ $? -eq 0 ]; then
        echo "success"
    else
        echo "failed"
        exit 1
    fi
}
