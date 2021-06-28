APPNAME := authorizer
VERSION := 1.0

build:
	mkdir -p build
	GOOS=$(GOOS) GOARCH=$(GOARCH) APPNAME=$(APPNAME) ./scripts/build.sh

run:
	./scripts/simple-run.sh

unit-test:
	./scripts/unit-test.sh

integration-test:
	./scripts/integration-test.sh

lint:
	./scripts/lint.sh

clean:
	rm -rf build
	rm -rf reports
	rm -rf ${APPNAME}.log

docker-build:
	./scripts/docker-build.sh

update-dependencies:
	./scripts/update-dependencies.sh


.PHONY: build run unit-test integration-test lint clean docker-build update-dependencies