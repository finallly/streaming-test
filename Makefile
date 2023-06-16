all: protoc client server

protoc:
	@echo "Generating Go files"
	cd src/proto && protoc --go_out=:. *.proto

server: protoc
	@echo "Building server"
	go build -o server \
		github.com/finallly/streaming-test/src/server

client: protoc
	@echo "Building client"
	go build -o client \
		github.com/finallly/streaming-test/src/client

.PHONY: client server protoc