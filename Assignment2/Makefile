
SERVER_HOST?=localhost
SERVER_PORT?=50000
MONGO_HOST?=localhost
MONGO_PORT?=27017
MONGO_USERNAME?=admin
MONGO_PASSWORD?=admin
MONGO_DB?=marketplace
POSTGRES_HOST?=localhost
POSTGRES_PORT?=5432
POSTGRES_USERNAME?=admin
POSTGRES_PASSWORD?=admin
POSTGRES_DB?=marketplace
POSTGRES_MAX_CONN?=500


build-server-seller:
	go build -o ./cmd/server/server-seller ./cmd/server

build-server-buyer:
	go build -o ./cmd/server/server-buyer ./cmd/server

build-client-seller:
	go build -o ./cmd/seller/client-seller ./cmd/seller

build-client-buyer:
	go build -o ./cmd/buyer/client-buyer ./cmd/buyer

build-test-latency:
	go build -o ./cmd/test_latency/test-latency ./cmd/test_latency

build-test-throughput:
	go build -o ./cmd/test_throughput/test-throughput ./cmd/test_throughput

build-all: build-server-seller build-server-buyer build-client-seller build-client-buyer build-test-latency build-test-throughput

run-postgres:
	docker compose -f ./deployment/docker/docker-compose.yaml up -d postgres pgbouncer
	docker compose -f ./deployment/docker/docker-compose.yaml ps

run-mongo:
	docker compose -f ./deployment/docker/docker-compose.yaml up -d mongodb
	docker compose -f ./deployment/docker/docker-compose.yaml ps

run-server-seller: build-server-seller
	./cmd/server/server-seller

run-server-buyer: build-server-buyer
	./cmd/server/server-buyer

run-client-seller: build-client-seller
	./cmd/seller/client-seller

run-client-buyer: build-client-buyer
	./cmd/buyer/client-buyer

run-test-latency: build-test-latency
	./cmd/test_latency/test-latency

run-test-throughput: build-test-throughput
	./cmd/test_throughput/test-throughput

build-proto:
	 protoc \
 	--proto_path=$(GOPATH)/src/github.com/adarshsrinivasan/DS_S24/library/proto\
 	--plugin=$(GOPATH)/bin/protoc-gen-go-grpc \
 	--go_out=$(GOPATH)/src \
 	--go-grpc_out=$(GOPATH)/src \
 	common.proto \
 	nosql-api.proto \
 	sql-api.proto
