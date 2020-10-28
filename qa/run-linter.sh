#!/bin/bash

# Script to check all Go scripts for linter issues

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"

IFS=$'\n'

# list of directories with sources to check
directories=$(cat ${SCRIPT_DIR}/directories.txt)

pass=0
fail=0

pushd "${SCRIPT_DIR}/.."


echo "----------------------------------------------------"
echo "Running Python linter against following directories:"
echo "$directories"
echo "----------------------------------------------------"
echo

getgolint="$(go get -u golang.org/x/lint/golint)"

# checks for the whole directories
for directory in $directories
do
    files=$(find "$directory" -name '*.go' -print)
    for source in $files
    do
        echo "Checking source $source"
        result=$(go run golang.org/x/lint/golint "$source")
        if [[ $result ]];
        then
            echo "$result"
            echo "    Fail"
            let "fail++"
        else
            echo "    Pass"
            let "pass++"
        fi
    done
done

popd

if [ $fail -eq 0 ]
then
    echo "All checks passed for $pass source files"
else
    let total=$pass+$fail
    echo "Linter fail, $fail source files out of $total source files need to be fixed"
    exit 1
fi
