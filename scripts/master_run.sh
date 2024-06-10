#!/bin/bash

# This script is used to run the master node of the cluster

timestamp=$(date +%Y-%m-%d-%H-%M-%S)
log_path="/opt/minik8s/logs/$timestamp"
mkdir -p $log_path

declare -A components_args
components_args["apiserver"]=""
components_args["scheduler"]=""
components_args["ctlmgr"]="127.0.0.1:8080"
components_args["sl_gtw"]=""
components_args["kubeproxy"]="192.168.1.12:8080"
components_args["jobserver"]=""

function killall(){
    echo "Killing all components..."
    for component in "${!components_args[@]}"; do
        echo "Killing $component..."
        pkill -f $component
    done
    echo "All components killed"
}

function start(){
    killall
    echo "Starting master node..."
    for component in "${!components_args[@]}"; do
        echo "Starting $component..."
        nohup ./$component ${components_args[$component]} > $log_path/$component.log 2>&1 &
    done
    echo "Master node started"
}

# 定义主程序，如果第一个参数为start，则执行start函数，否则执行killall函数
# 无参数执行start

if [ $# -eq 0 ]; then
    start
elif [ "$1" == "start" ]; then
    start
elif [ "$1" == "killall" ]; then
    killall
else 
    echo "Usage: $0 [start|killall]"
    exit 1
fi

exit 0
