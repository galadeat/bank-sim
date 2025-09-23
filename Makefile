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

pb: 
	protoc --go_out=paths=source_relative:. \
       --go-grpc_out=paths=source_relative:. \
       api/proto/account/v1/account.proto

