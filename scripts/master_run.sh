#!/bin/bash

# This script is used to run the master node of the cluster
echo "Starting master node..."
./apiserver &
./scheduler &
./ctlmgr &
echo "Master node started"
