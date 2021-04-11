#!/bin/bash -ex
set -exv

mkdir -p $HOME/.crda

export AUTH_TOKEN=${THREE_SCALE_KEY};
export CRDA_KEY=${CRDA_KEY};
export HOST=https://f8a-analytics-preview-2445582058137.staging.gw.apicast.io
export snyk_token=${SNYK_TOKEN}

printf 'auth_token: %s\ncrda_key: %s\nhost: %s' "${AUTH_TOKEN}" "${CRDA_KEY}" "${HOST}" >> $HOME/.crda/config.yaml

cat $HOME/.crda/config.yaml

go mod tidy 

go test -test.v -ginkgo.failFast -ginkgo.focus="PR ACCEPTANCE TESTS"