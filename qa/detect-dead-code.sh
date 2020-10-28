#!/bin/bash

# Script to check all Python scripts for PEP-8 issues

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"

IFS=$'\n'

# list of directories with sources to check
directories=$(cat ${SCRIPT_DIR}/directories.txt)

pass=0
fail=0

pushd "${SCRIPT_DIR}/.."

getgocheck="$(go get -u github.com/opennota/check)"

# run the vulture for all files that are provided in $1
function check_files() {
    for source in $1
    do
        echo "Checking $source"
        result=$(go run github.com/opennota/check/cmd/varcheck "$source")
        if [ $? -eq 0 ]
        then
            echo "    Pass"
            let "pass++"
        elif [ $? -eq 2 ]
        then
            echo "    Illegal usage (should not happen)"
            exit 2
        else
            echo "$result"
            echo "    Fail"
            let "fail++"
        fi
    done
}


echo "----------------------------------------------------"
echo "Checking source files for dead code and unused imports"
echo "in following directories:"
echo "$directories"
echo "----------------------------------------------------"
echo

# checks for the whole directories
for directory in $directories
do
    files=$(find "$directory" -name '*.go' -print)

    check_files "$files"
done

popd

if [ $fail -eq 0 ]
then
    echo "All checks passed for $pass source files"
else
    let total=$pass+$fail
    echo "$fail source files out of $total files seems to contain dead code"
    exit 1
fi

