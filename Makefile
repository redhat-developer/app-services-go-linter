binary:
	go build -o go-i18n-linter ./cmd/go-i18n-linter

lint:
	golangci-lint run cmd/... pkg/...

test:
	go test ./pkg/...