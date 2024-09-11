#!/bin/bash

for node in $(kubectl get nodes --selector='!node-role.kubernetes.io/control-plane' -o name); do
  kubectl label $node node-role.kubernetes.io/worker=""
done
