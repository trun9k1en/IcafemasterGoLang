.PHONY: build run test clean deps docker-up docker-down

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod tidy
	go mod download

# Start MongoDB with Docker
docker-up:
	docker-compose up -d

# Stop MongoDB Docker
docker-down:
	docker-compose down

# Build Docker image
docker-build:
	docker build -t icafe-registration .

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Generate swagger docs (if using swaggo)
swagger:
	swag init -g cmd/api/main.go -o docs
