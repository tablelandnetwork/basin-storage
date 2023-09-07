# Lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2 run
.PHONY: lint

# Test
test: 
	go test ./... -short -race
.PHONY: test

mocks: clean-mocks
	go run github.com/vektra/mockery/v2@v2.14.0 --name=GCS --recursive --with-expecter
.PHONY: mocks

clean-mocks:
	rm -rf mocks
.PHONY: clean-mocks

uploader-local:
	WEB3STORAGE_TOKEN=$(W3S_TOKEN) CRDB_CONN_STRING=$(CRDB_CONN) FUNCTION_TARGET=Uploader go run cmd/main.go
.PHONY: uploader-local

uploader-deploy:
	gcloud functions deploy go-finalize-function \
	--gen2 \
	--runtime=go120 \
	--region=us-central1 \
	--source=. \
	--entry-point=Uploader \
	--trigger-event-filters="type=google.cloud.storage.object.v1.finalized" \
	--trigger-event-filters="bucket=tableland-entrypoint"  \
	--memory 8192MB \
	--timeout 3600s \
	--set-env-vars WEB3STORAGE_TOKEN=$(W3S_TOKEN),CRDB_CONN_STRING="$(CRDB_CONN)"
.PHONY: uploader-deploy

statuschecker-local:
	WEB3STORAGE_TOKEN=$(W3S_TOKEN) CRDB_CONN_STRING=$(CRDB_CONN) FUNCTION_TARGET=StatusChecker go run cmd/main.go
.PHONY: statuschecker-local

ethereum:
	go run github.com/ethereum/go-ethereum/cmd/abigen@v1.11.4 --abi ./evm/basin_storage/out/BasinStorage.sol/BasinStorage.abi.json --bin ./evm/basin_storage/out/BasinStorage.sol/BasinStorage.bin --pkg ethereum --type Contract --out pkg/ethereum/contract.go
.PHONY: ethereum	