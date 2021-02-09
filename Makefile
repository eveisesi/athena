gqlgen:
	go run github.com/99designs/gqlgen generate --config .config/gqlgen.yml

run-server: build
	.build/athena server

run-processor: build
	.build/athena processor

build:
	go build -o .build/athena cmd/app/*.go