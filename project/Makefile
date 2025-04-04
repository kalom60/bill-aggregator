BROKER_BINARY=brokerApp
USER_BINARY=userApp
PROVIDER_BINARY=providerApp
ACCOUNT_BINARY=accountApp
MOCK_BINARY=mockApp
BILL_BINARY=billApp
LISTENER_BINARY=listenerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
# up_build: build_broker build_user build_provider build_account build_mock build_bill
up_build: build_broker build_user build_provider build_account build_mock build_bill build_listener
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_user: builds the user binary as a linux executable
build_user:
	@echo "Building user binary..."
	cd ../user-service && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_BINARY} ./cmd/api
	@echo "Done!"

## build_provider: builds the provider binary as a linux executable
build_provider:
	@echo "Building provider binary..."
	cd ../utility-provider-service && env GOOS=linux CGO_ENABLED=0 go build -o ${PROVIDER_BINARY} ./cmd/api
	@echo "Done!"

## build_account: builds the account linking binary as a linux executable
build_account:
	@echo "Building account binary..."
	cd ../account-linking-service && env GOOS=linux CGO_ENABLED=0 go build -o ${ACCOUNT_BINARY} ./cmd/api
	@echo "Done!"

## build_mock: builds the mock api binary as a linux executable
build_mock:
	@echo "Building mock binary..."
	cd ../mock-api-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MOCK_BINARY} ./cmd/api
	@echo "Done!"

## build_mock: builds the bill aggregation binary as a linux executable
build_bill:
	@echo "Building bill binary..."
	cd ../bill-aggregation-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BILL_BINARY} ./cmd/api
	@echo "Done!"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building listener binary..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Done!"

