#!/usr/bin/env bash

# List of tools that used
tools=(
	gotest.tools/gotestsum
)

# Install missed tools
for tool in ${tools[@]}; do
	which $(basename ${tool}) > /dev/null || go get -u -v ${tool}
done

echo "Running integration tests."

# Generate tests report
gotestsum -- -tags=integration ./cmd/authorizer; test ${PIPESTATUS[0]} -eq 0 || status=${PIPESTATUS[0]}

exit ${status:-0}