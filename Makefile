.PHONY: run-server run-client quickstart test pb

run-server:
	go run ./server


run-client:
	go run ./client

quickstart:
	@go run ./server & \
	SERVER_PID=$$!; \
	echo "Waiting for gRPC server on port 50051"; \
	until nc -z localhost 50051; do sleep 1; done; \
	echo "server is up"; \
	go run ./client; \
	kill $$SERVER_PID

test:
	go test ./tests/... -v

PROTOC = protoc
PROTO_DIR = api/proto
PROTO_FILES = $(shell find $(PROTO_DIR) -name "*.proto")

pb:
	rm -rf api/proto/**/*.pb.go
	$(PROTOC) \
		--proto_path=$(PROTO_DIR) \
    	--go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
    	--go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
    	$(PROTO_FILES)


