# Lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2 run
.PHONY: lint

# Test
test: 
	go test ./... -short -race
.PHONY: test

uploader-test:
	go test ./... -timeout 30s -run ^TestUploader$
.PHONY: uploader-test	

checker-test:
	go test ./... -timeout 5m -run ^TestChecker$
.PHONY: checker-test

mocks: clean-mocks
	go run github.com/vektra/mockery/v2@v2.14.0 --name=GCS --recursive --with-expecter
.PHONY: mocks

clean-mocks:
	rm -rf mocks
.PHONY: clean-mocks

uploader-local:
	FUNCTION_TARGET=Uploader go run cmd/main.go
.PHONY: uploader-local

uploader-deploy:
	gcloud functions deploy go-finalize-function \
	--gen2 \
	--runtime=go120 \
	--region=us-east1 \
	--source=. \
	--entry-point=Uploader \
	--trigger-event-filters="type=google.cloud.storage.object.v1.finalized" \
	--trigger-event-filters="bucket=tableland-basin-staging"  \
	--memory 32768MB \
	--timeout 3600s \
	--env-vars-file uploader.env.yml
.PHONY: uploader-deploy

checker-local:
	FUNCTION_TARGET=StatusChecker go run cmd/main.go
.PHONY: checker-local

checker-deploy:
	gcloud functions deploy go-http-function \
  	--gen2 \
	--region=us-east1 \
	--runtime=go120 \
	--source=. \
	--entry-point=StatusChecker \
	--trigger-http \
	--memory 8192MB \
	--timeout 600s \
	--run-service-account basin-status-checker-gcf@textile-310716.iam.gserviceaccount.com \
	--env-vars-file checker.env.yml
.PHONY: checker-deploy

ethereum:
	go run github.com/ethereum/go-ethereum/cmd/abigen@v1.12.2 --abi ./evm/basin_storage/out/BasinStorage.sol/BasinStorage.abi.json --bin ./evm/basin_storage/out/BasinStorage.sol/BasinStorage.bin --pkg ethereum --type Contract --out pkg/ethereum/contract.go
.PHONY: ethereum	