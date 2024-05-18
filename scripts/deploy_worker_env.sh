#!/bin/bash

. ./scripts/deploy_common_env.sh

worker_ip=("worker-1" "worker-2")

declare -A worker_identity_file
worker_identity_file["worker-1"]="~/.ssh/k8s-default.pem"
worker_identity_file["worker-2"]="~/.ssh/k8s-default.pem"

declare -A worker_user
worker_user["worker-1"]="root"
worker_user["worker-2"]="root"
