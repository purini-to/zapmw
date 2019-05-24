.PHONY: test fmt lint

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...

fmt:
	go fmt ./...

lint:
	golangci-lint run --exclude-use-default=false ./...
