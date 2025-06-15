server: build
	@go run . -mode=server


client: build
	@go run . -mode=client

build:
	@go build -o bin/mq .
