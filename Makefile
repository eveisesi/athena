gqlgen:
	go run github.com/99designs/gqlgen generate --config .config/gqlgen.yml

run-server: build
	./athena server

run-processor: build
	./athena processor

build:
	go build -o athena cmd/app/*.go