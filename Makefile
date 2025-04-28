include .env
MIGRATIONS_DIR=./internal/infra/persistence/postgres/migrations

run:
	@go run ./cmd/main.go
test:
	@go test -v ./...
migration:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(filter-out $@, $(MAKECMDGOALS))
migrate-up:
	@migrate -database $(DATABASE_ADDR) -path $(MIGRATIONS_DIR) up
migrate-down:
	@migrate -database $(DATABASE_ADDR) -path $(MIGRATIONS_DIR) down
