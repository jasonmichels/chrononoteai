# Variables
APP_NAME = chrononoteai
SRC = ./main.go
TEST_DIR = ./...

# Default target
all: build

# Build the application
build:
	@echo "Building the application..."
	@go build -o $(APP_NAME) $(SRC)

# Run tests
test:
	@echo "Running tests..."
	@go test $(TEST_DIR) -v

# Clean up
clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)

# Rebuild
rebuild: clean build

# Run
run: $(APP_NAME)

.PHONY: all build test clean rebuild
