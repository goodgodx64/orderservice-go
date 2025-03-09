#Makefile
all-srv: protos build-srv run-srv
protos:
	@echo "Generating proto files..."
	@protoc --go_out=./pkg/api/test --go_opt=paths=source_relative  --go-grpc_out=./pkg/api/test --go-grpc_opt=paths=source_relative order.proto
build-srv:
	@echo "Building server..."
	@mkdir bin
	@go build -o bin/server cmd/main.go
run-srv: build-srv
	@echo "Running server..."
	@./bin/server