SERVER_HOST?=localhost
SERVER_PORT?=50000
MONGO_HOST?=localhost
MONGO_PORT?=50001
MONGO_USER?=admin
MONGO_PASSWORD?=admin
MONGO_DB?=marketplace
POSTGRES_HOST?=localhost
POSTGRES_PORT?=50002
POSTGRES_USER?=admin
POSTGRES_PASSWORD?=admin
POSTGRES_DB?=marketplace
POSTGRES_MAX_CONN?=500

# use the rest as arguments
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
# ...and turn them into do-nothing targets
$(eval $(RUN_ARGS):;@:)


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

run-sql-server:
	docker compose -d postgres pgbouncer
	docker compose ps

run-nosql-server:
	docker compose up -d mongodb
	docker compose ps

run-server-seller: build-server-seller
	./cmd/server/server-seller

run-server-buyer: build-server-buyer
	./cmd/server/server-buyer

run-client-seller: build-client-seller
	./cmd/seller/client-seller

run-client-buyer: build-client-buyer
	./cmd/buyer/client-buyer

run-test-latency: build-test-latency
	./cmd/test_latency/test-latency $(RUN_ARGS)

run-test-throughput: build-test-throughput
	./cmd/test_throughput/test-throughput $(RUN_ARGS)
 	
clean:
	docker compose down -v || true
	rm -rf ./cmd/server/server-seller || true
	rm -rf ./cmd/server/server-buyer || true
	rm -rf ./cmd/seller/client-seller || true
	rm -rf ./cmd/buyer/client-buyer || true
	rm -rf ./cmd/test_latency/test-latency || true
	rm -rf ./cmd/test_throughput/test-throughput || true 
