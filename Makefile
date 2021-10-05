binary:
	go build -o app-services-go-linter ./cmd/app-services-go-linter
.PHONY: binary

plg:
	go build -buildmode=plugin plugin/appservicesgolinter.go
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