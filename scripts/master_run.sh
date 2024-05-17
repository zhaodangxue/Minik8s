#!/bin/bash

# This script is used to run the master node of the cluster
echo "Starting master node..."
./apiserver &
./scheduler &
./ctlmgr 127.0.0.1:8080 &
echo "Master node started"
