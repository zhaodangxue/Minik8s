#!/bin/bash

. ./scripts/deploy_worker_env.sh

ssh_execute_command() {
    local command=$1
    local description=$2

    echo "Executing: $description"
    ssh -i $deploy_identity_file $deploy_user@$deploy_url $command
    if [ $? -eq 0 ]; then
        echo "success"
    else
        echo "failed"
        exit 1
    fi
}

echo "Deploying to worker..."
for i in ${!worker_ip[@]}; do
    echo "Deploying to worker $i: ${worker_ip[$i]}"
    deploy_url=${worker_ip[$i]}
    deploy_identity_file=${worker_identity_file[$deploy_url]}
    deploy_user=${worker_user[$deploy_url]}
    ssh_execute_command "echo Connected to worker" "connect to worker"
    ssh_execute_command "mkdir -p $deploy_dir" "create worker deploy directory"
    ssh_execute_command "rm -rf $deploy_dir/*" "clean worker deploy directory"
    execute_command "scp -i $deploy_identity_file -r $build_dir/* $deploy_user@$deploy_url:$deploy_dir" "copy product to deploy directory"
    ssh_execute_command "chmod +x $deploy_dir/*" "set execution permissions"
done
