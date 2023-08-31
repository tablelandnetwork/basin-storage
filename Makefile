# Lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2 run
.PHONY: lint

# Test
test: 
	go test ./... -short -race
.PHONY: test

mocks: clean-mocks
	go run github.com/vektra/mockery/v2@v2.14.0 --name=GCSOps --recursive --with-expecter
.PHONY: mocks

clean-mocks:
	rm -rf mocks
.PHONY: clean-mocks