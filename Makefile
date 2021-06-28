APPNAME := authorizer
VERSION := 1.0

## build: Compile project.
build:
	mkdir -p build
	GOOS=$(GOOS) GOARCH=$(GOARCH) APPNAME=$(APPNAME) ./scripts/build.sh

## run: Build and execute project using testdata/operations.
run:
	./scripts/simple-run.sh

## unit-test: Execute unit tests.
unit-test:
	./scripts/unit-test.sh

## integration-test: Execute integration tests.
integration-test:
	./scripts/integration-test.sh

## lint: Execute linter using the rules in .golangci.yml.
lint:
	./scripts/lint.sh

## clean: Delete build directory and log.
clean:
	rm -rf build
	rm -rf ${APPNAME}.log
	rm -rf cover.out

## docker-build: Build project and create container with it.
docker-build:
	./scripts/docker-build.sh

## docker-run: Execute the application within a container with testdata.
docker-run:
	./scripts/docker-run.sh

## update-dependencies: Update all golang dependencies.
update-dependencies:
	./scripts/update-dependencies.sh


.PHONY: build run unit-test integration-test lint clean docker-build update-dependencies docker-run

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo