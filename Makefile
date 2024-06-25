build:
	@go build -o bin/airbnbreplica cmd/main.go

test:
	@go test -v ./...

run : build
	@./bin/airbnbreplica