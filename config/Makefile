.PHONY: build run swag clean

build:
	go build -o config-api cmd/main.go

run:
	go run cmd/main.go

swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./internal/swagger_expose.go -o ./api

clean:
	rm -f config-api
	rm -rf api
