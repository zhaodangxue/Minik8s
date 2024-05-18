#!/bin/bash

. ./scripts/deploy_master_env.sh

echo "Deploying to master..."

ssh_execute_command() {
    local command=$1
    local description=$2

    echo "Executing: $description"
    ssh -i $master_identity_file $master_user@$master_url $command
    if [ $? -eq 0 ]; then
        echo "$description executed successfully"
    else
        echo "Failed to execute $description"
        exit 1
    fi
}

# Check directory
if [ ! -d $build_dir ]; then
    echo "Build directory not found"
    exit 1
fi

ssh_execute_command "echo Connected to master" "connect to master"

# Copy product to master
ssh_execute_command "mkdir -p $deploy_dir" "create master deploy directory"
ssh_execute_command "rm -rf $deploy_dir/*" "clean master deploy directory"
execute_command "scp -i $master_identity_file -r $build_dir/* $master_user@$master_url:$deploy_dir" "copy product to deploy directory"
ssh_execute_command "chmod +x $deploy_dir/*" "set execution permissions"
