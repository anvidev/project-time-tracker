GOOSE_DRIVER = turso
GOOSE_DBSTRING = $(TURSO_DATABASE_URL)?authToken=$(TURSO_AUTH_TOKEN)
GOOSE_MIGRATION_DIR = ./cmd/migrate/migrations

build:
	@echo "Building executable..."
	@go build -o ./bin/api ./cmd/api/*.go

serve:
	@echo "Running executable..."
	@./bin/api

run: build serve

migrate:
	@read -p "Enter the sequence name: " SEQ; \
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) goose create $${SEQ} sql

migrate-up:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) goose up

migrate-down:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) goose down

migrate-status:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) goose status
