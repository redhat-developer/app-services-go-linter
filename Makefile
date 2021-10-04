binary:
	go build -o go-i18n-linter ./cmd/go-i18n-linter
.PHONY: binary

plg:
	go build -buildmode=plugin plugin/goi18nlinter.go
.PHONY: plg

lint:
	golangci-lint run cmd/... pkg/...
.PHONY: lint

test:
	go test ./pkg/...
.PHONY: test

format:
	@go mod tidy

	@gofmt -w `find . -type f -name '*.go'`
.PHONY: format