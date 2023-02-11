#!/bin/bash -e

help()
{
   # Display Help
   echo "Usage: ./run.sh [ -h | -help | -f | -u | -d ] or none"
   echo
   echo "-help           See this message."
   echo "-f              Run functional tests."
   echo "-u              Run unit tests."
   echo "-d              Run app with the default params."
   echo "-c              Clean log dir."
   echo "nothing         Redirect to the go executable."
}

run_func_tests() {
    go run cmd/func_tests/func_tests.go
}

run_unit_tests() {
    go test -v ./...
}

run_app_default() {
    go run cmd/listam_parser/*.go -p 200000
}

run_app() {
    go run cmd/listam_parser/*.go "$@"
}

clean_logdir() {
    rm -r log/*
}

##############
### Main #####
##############

# Doesn`t work wtf ?`
# cd "$(dirname "$(readlink -f "$0")")"
# cd "$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

if  [ "$1" == "-help" ]; then
    help
elif [ "$1" == "-f" ]; then
    run_func_tests
elif [ "$1" == "-u" ]; then
    run_unit_tests
elif [ "$1" == "-d" ]; then
    run_app_default
elif [ "$1" == "-c" ]; then
    clean_logdir
else 
    run_app "$@"
fi
