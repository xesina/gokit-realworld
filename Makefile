#@IgnoreInspection BashAddShebang
export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export DEBUG=true
export APP=gokit-realworld
export LDFLAGS="-w -s"

all: build test

build:
	go build -race  .

build-static:
	CGO_ENABLED=0 go build -race -v -o $(APP) -a -installsuffix cgo -ldflags $(LDFLAGS) .

run:
	go run -race .

############################################################
# Test
############################################################

test:
	go test -v -race ./...

container:
	docker build -t echo-realworld .

run-container:
	docker run --rm -it echo-realworld

.PHONY: build run build-static test container

############################################################
# Lint
############################################################

.PHONY: install-linter check-linter lint
install-linter:
	cd /tmp && GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0 && cd -

check-linter:
ifeq ($(shell which golangci-lint), )
	make install-linter
endif

lint: check-linter
	go mod download
	golangci-lint run --fast

.PHONY: fmt
fmt:
	go fmt ./...
