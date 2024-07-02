build:
	@go build -o bin/airbnbreplica cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/airbnbreplica

start-redis:
	@docker start airbnb-redis

stop-redis:
	@docker stop airbnb-redis

.PHONY: build test run start-redis stop-redis
