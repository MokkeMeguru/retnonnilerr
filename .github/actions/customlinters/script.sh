#!/bin/bash
TARGET=customlinters

go build -o bin/customlinters tools/analysis/main.go
ANALYSIS_RESULT=$(./bin/customlinters ./... 2>&1)
ANALYSIS_RESULT_EXIT_CODE=$?

set -eo pipefail

if [ ${ANALYSIS_RESULT_EXIT_CODE} -eq 3 ]; then
    echo "${ANALYSIS_RESULT}" | grep -E -v "^#" |
        reviewdog \
            -name="${TARGET}" \
            -f="golint" \
            -reporter="github-pr-review" \
            -filter-mode="nofilter" \
            -fail-on-error="false" \
            -level="warning"
elif [ ${ANALYSIS_RESULT_EXIT_CODE} -eq 0 ]; then
    exit 0
elif [ ${PR_NUMBER} -ne "" ]; then
    gh pr comment "${PR_NUMBER}" --body "failed to run custom linters at [the action](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID})"
fi
