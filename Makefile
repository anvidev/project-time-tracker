build:
	@echo "Building executable..."
	@go build -o ./bin/api ./cmd/api

serve:
	@echo "Running executable"
	@./bin/api

run: build serve
