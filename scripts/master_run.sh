#!/bin/bash

# This script is used to run the master node of the cluster
echo "Starting master node..."
nohup ./apiserver &
nohup ./scheduler &
nohup ./ctlmgr 127.0.0.1:8080 &
echo "Master node started"

exit 0
