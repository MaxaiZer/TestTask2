BINARY_NAME=app
ifeq ($(OS),Windows_NT)
    BINARY_NAME := $(BINARY_NAME).exe
endif

test:
	@echo "Running tests..."
	go test -v ./...

swag:
	@echo "Generating Swagger documentation..."
	swag init -g ./cmd/main.go

build:
	@echo "Building the application..."
	go build -o $(BINARY_NAME) ./cmd

run: build
	@echo "Running the application..."
	./$(BINARY_NAME)

clean:
	@echo "Cleaning the build..."
	rm -f $(BINARY_NAME)

.PHONY: build clean run test swag