#!/bin/bash

SERVER_UUID="$1"

function create_new_snapshot {
    local SERVER_UUID="$1"
    echo -n "Creating new snapshot for server ${SERVER_UUID}."
    local SNAP_CREATION=$(titan-sc snap create -u ${SERVER_UUID} | \
        jq -r '.|if .error then .code else "SUCCESS" end')
    if [ "${SNAP_CREATION}" == "SUCCESS" ]; then
        while [ $(titan-sc snap ls -u ${SERVER_UUID} | \
            jq -r '.[]|select(.state=="creating")|.uuid' | wc -l) -gt 0 ]; do
            echo -n "."
            sleep 2
        done
        echo -e "\nNew snapshot created."
        return 0
    elif [ "${SNAP_CREATION}" == "SNAPSHOT_CREATE_FAIL_LIMIT_EXCEEDED" ]; then
        echo -e "\nSnapshot limit reached for server ${SERVER_UUID}."
        return 1
    else
        echo -e "\nUnable to create new snapshot: ${SNAP_CREATION}."
        exit -1
    fi
}

function delete_oldest_snapshot {
    local SERVER_UUID="$1"
    local OLDEST_SNAP_UUID=$(titan-sc snap ls -u ${SERVER_UUID} | \
        jq -r '.|if type!="array" then "ERROR" else .[].created_at|=(.|split(".")[0]|
               strptime("%Y-%m-%dT%H:%M:%S")|mktime)|sort_by(.created_at)[0].uuid end')
    if [ "${OLDEST_SNAP_UUID}" == "ERROR" ]; then
        echo "Unable to retrieve snapshots of server ${SERVER_UUID}."
        exit -1
    fi
    echo -n "Deleting oldest snapshot ${OLDEST_SNAP_UUID}."
    local SNAP_DELETION=$(titan-sc snap del -u ${SERVER_UUID} -s ${OLDEST_SNAP_UUID} | \
        jq -r '.|if .error then .code else "SUCCESS" end')
    if [ "${SNAP_DELETION}" != "SUCCESS" ]; then
        echo -e "\nUnable to delete snapshot ${OLDEST_SNAP_UUID}: ${SNAP_DELETION}"
        exit -1
    fi
    while [ $(titan-sc snap ls -u ${SERVER_UUID} | \
        jq -r '.[]|select(.state=="deleting")|.uuid' | wc -l) -gt 0 ]; do
        echo -n "."
        sleep 2
    done
    echo -e "\nSnapshot deleted."
}

# Create a new snapshot (delete the oldest one when limit has been reached)
create_new_snapshot ${SERVER_UUID}
if [ $? -eq 1 ]; then
    delete_oldest_snapshot ${SERVER_UUID}
    create_new_snapshot ${SERVER_UUID}
fi
exit 0
