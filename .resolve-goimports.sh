#!/bin/bash
COMMAND="goimports-reviser"
ORDER="std,project,general,company"

install() {
    if ! command -v "${COMMAND}" >/dev/null 2>&1; then
        go install -v github.com/incu6us/goimports-reviser/v3@v3.4.1
    fi
}

general_format() {
    "${COMMAND}" -rm-unused -set-alias -imports-order "${ORDER}" -format -recursive "$(pwd)"
}

quick_format() {
    "${COMMAND}" -imports-order "${ORDER}" -format -recursive "$(pwd)"
}

case "$1" in
"-q")
    quick_format
    ;;
"-g")
    general_format
    ;;
"-i")
    install
    ;;
*)
    echo "Usage: $0 {-q|-g|-i}"
    exit 2
    ;;
esac

exit $?
