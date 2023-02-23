#!/bin/bash -e

help()
{
   # Display Help
   echo "Usage: ./run.sh [ ftests | utests | default | cleanlog | params ]"
   echo
   echo "ftests             Run functional tests."
   echo "utests             Run unit tests."
   echo "defaut             Run app with the default params."
   echo "cleanlog           Clean log dir."
   echo "params             Redirect to the go executable."
}

run_func_tests() {
    go run cmd/func_tests/func_tests.go
}

run_unit_tests() {
    go test -v ./internal/...
}

run_app_default() {
    go run cmd/listam_parser/*.go -p 200000 -t 0
}

run_app() {
    go run cmd/listam_parser/*.go "$@"
}

clean_logdir() {
    rm -rf log/*.json && rm -rf log/*.log
}

##############
### Main #####
##############

# Doesn`t work wtf ?`
# cd "$(dirname "$(readlink -f "$0")")"
# cd "$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

if [ "$1" == "ftests" ]; then
    run_func_tests
elif [ "$1" == "utests" ]; then
    run_unit_tests
elif [ "$1" == "default" ]; then
    run_app_default
elif [ "$1" == "cleanlog" ]; then
    clean_logdir
elif [ "$1" == "params" ]; then
    shift
    run_app "$@"
else 
    help    
fi
