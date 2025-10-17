.PHONY: build run-server run-client quickstart tests pb mock

run-server:
	go run ./server


run-client:
	go run ./client

build:
	@go build -o bin/server ./cmd/server

quickstart: build
	@./bin/server & \
	SERVER_PID=$$!; \
	echo "Waiting for gRPC servers..."; \
	until nc -z localhost 50051 && nc -z localhost 50052; do sleep 1; done; \
	echo "servers are up"; \
	go run ./cmd/client; \
	sleep 1; \
	kill $$SERVER_PID; \


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


mock:
	mockgen \
      -destination=mocks/mock_user_client.go \
      -package=mocks \
      -source=api/proto/user/v1/user_grpc.pb.go


tests:
	go test -coverprofile=coverage.out ./client/... ./server/... && go tool cover -func coverage.out