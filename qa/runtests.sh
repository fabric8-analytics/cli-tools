#!/bin/bash

# Script to check all Go scripts for linter issues

SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"

IFS=$'\n'

pass=0
fail=0

pushd "${SCRIPT_DIR}/.."

result=$(go test -v -cover ./...)
passFail=$(echo $result | grep "FAIL" | wc -l)
if [[ $passFail -eq 0 ]];
then
    echo "    Pass"
    let "pass++"
else
    echo "$result"
    echo "    Fail"
    let "fail++"
fi

popd

if [ $fail -eq 0 ]
then
    echo "All tests passed"
else
    let total=$pass+$fail
    echo "Test case failed"
    exit 1
fi

