#!/bin/bash -ex
set -exv

mkdir -p $HOME/.crda
 
export CRDA_AUTH_TOKEN=${THREE_SCALE_KEY:-3e42fa66f65124e6b1266a23431e3d08};
export CRDA_KEY=${CRDA_KEY:-d931dd95-ab1f-4f74-9a9f-fb50f60e4ea9};
export CRDA_HOST=https://f8a-analytics-preview-2445582058137.staging.gw.apicast.io
export snyk_token=${SNYK_TOKEN}

printf 'crda_auth_token: %s\nconsent_telemetry: false\ncrda_key: %s\nhost: %s' "${CRDA_AUTH_TOKEN}" "${CRDA_KEY}" "${HOST}" >> $HOME/.crda/config.yaml

cat $HOME/.crda/config.yaml

go test -test.v -ginkgo.failFast -ginkgo.reportFile ginkgo.report -ginkgo.focus="PR ACCEPTANCE TESTS" 