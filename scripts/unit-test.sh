#!/usr/bin/env bash
set -e

mkdir -p reports

COVER_FILE="reports/cover.out"
COVERAGE_REPORT="reports/coverage.xml"
JUNIT_REPORT="reports/junit-report.xml"

# List of tools that used to generate Quality Gate reports
tools=(
	github.com/axw/gocov/gocov
	github.com/AlekSi/gocov-xml
	gotest.tools/gotestsum
)

# Install missed tools
for tool in ${tools[@]}; do
	which $(basename ${tool}) > /dev/null || go get -u -v ${tool}
done

echo "Running unit tests."

# Generate tests report
gotestsum --junitfile "${JUNIT_REPORT}" -- -coverprofile=${COVER_FILE} ./...; test ${PIPESTATUS[0]} -eq 0 || status=${PIPESTATUS[0]}

# Print code coverage details
go tool cover -func "${COVER_FILE}"

# Generate coverage report
echo "Generate coverage report."
gocov convert "${COVER_FILE}" | gocov-xml  > ${COVERAGE_REPORT}; test ${PIPESTATUS[0]} -eq 0 || status=${PIPESTATUS[0]}

exit ${status:-0}