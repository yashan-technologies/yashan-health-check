#!/bin/bash
USER=$(echo "$CI_COMMIT_AUTHOR" | cut -d ' ' -f 1)
ARTIFACTS_DIR=$(date -d "${CI_JOB_STARTED_AT}" +"%Y-%m-%d-%H-%M-%S")-MR-${CI_MERGE_REQUEST_IID}
ARTIFACTS_PATH=/home/gitlab/artifacts/$USER/${ARTIFACTS_DIR}

check_ret() {
    ret=$1
    name=$2
    if [ "$ret" -ne 0 ]; then
        echo "[${name}] failed, ret=${ret}"
        exit "$ret"
    fi
    return
}

dump() {
    mkdir -p "$ARTIFACTS_PATH"
    check_ret $? "make dir"
    if [ -f "./code_check.txt" ]; then
        cp -fr ./code_check.txt "$ARTIFACTS_PATH"
    fi
    check_ret $? "copy code_check.txt"

    if [ -d "./build" ]; then
        cp -fr ./build "$ARTIFACTS_PATH"
    fi
    check_ret $? "copy build"

    if [ -d "./unittest" ]; then
        cp -fr ./unittest "$ARTIFACTS_PATH"
    fi
    check_ret $? "copy unittest"
}

show() {
    echo http://"$CI_RUNNER":8888/"$USER"/"${ARTIFACTS_DIR}"/
}

case "$1" in
dump)
    dump
    ;;
show)
    show
    ;;
*)
    echo "Usage: $0 {dump|show}"
    exit 2
    ;;
esac

exit $?
