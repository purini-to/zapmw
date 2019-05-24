.PHONY: test fmt lint

test:
	go test -cover -v ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run --exclude-use-default=false ./...
