#!/bin/bash

echo "Starting DFS Clustor..."

mkdir -p logs

echo "Starting NameNode..."
../namenode/bin/namenode > logs/namenode.log 2>&1 &

echo "Starting Datanode1..."
../bin/server1  > logs/server1.log 2>&1 &

echo "Starting Datanode2..."
../bin/server2  > logs/server2.log 2>&1 &

echo "Starting Datanode3..."
../bin/server3  > logs/server3.log 2>&1 &

echo "All Nodes started!"