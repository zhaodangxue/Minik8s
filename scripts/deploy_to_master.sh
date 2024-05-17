#!/bin/bash

. ./scripts/deploy_env.sh

master_url=$MASTER_URL
master_deploy_dir=$MASTER_DEPLOY_DIR
master_user=root
master_identity_file=$MASTER_IDENTITY_FILE

build_dir=./build

echo "Deploying to master..."

execute_command() {
    local command=$1
    local description=$2

    echo "Executing: $description"
    $command
    if [ $? -eq 0 ]; then
        echo "$description executed successfully"
    else
        echo "Failed to execute $description"
        exit 1
    fi
}

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
ssh_execute_command "mkdir -p $master_deploy_dir" "create master deploy directory"
ssh_execute_command "rm -rf $master_deploy_dir/*" "clean master deploy directory"
execute_command "scp -i $master_identity_file -r $build_dir/* $master_user@$master_url:$master_deploy_dir" "copy product to deploy directory"
ssh_execute_command "chmod +x $master_deploy_dir/*" "set execution permissions"

# Restart service
ssh -i $master_identity_file $master_user@$master_url "cd $master_deploy_dir && setsid -f ./master_run.sh > /dev/null"
echo "Service restarted"
