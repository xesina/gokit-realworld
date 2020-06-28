.PHONY: install-linter check-linter lint
install-linter:
	cd /tmp && GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0 && cd -

check-linter:
ifeq ($(shell which golangci-lint), )
	make install-linter
endif

lint: check-linter
	go mod download
	golangci-lint run -n

.PHONY: fmt
fmt:
	go fmt ./...
