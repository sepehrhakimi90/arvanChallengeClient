#!/usr/bin/env bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
echo "${SCRIPT_DIR}"
export SERVER_HOST=192.168.122.1
export SERVER_PORT=8080
LOG_FILE="${SCRIPT_DIR}/client.log"
${SCRIPT_DIR}/client >> ${LOG_FILE} 2>&1 &
disown