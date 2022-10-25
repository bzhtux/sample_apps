#!/usr/bin/env bash

GetPod(){
    SEARCH_NAME=$1
    kubectl get pods --no-headers=true | grep "${SEARCH_NAME}" | awk '{print $1}'
}

IsPodRunning() {
    POD_NAME=$1
    if [ "$(kubectl get pod "${POD_NAME}" -o jsonpath="{.status.containerStatuses[*].state.running}")" == "" ];then
        return 1
    else
        return 0
    fi
}

WaitForPod() {
    POD_NAME=$1
    SLEEP=$2
    while [ "$(kubectl get pod "${POD_NAME}" -o jsonpath="{.status.containerStatuses[*].state.running}")" == "" ]
do
    echo ">>> Waiting for ${POD_NAME} to be up and running ..."
    sleep "${SLEEP}"
done
}