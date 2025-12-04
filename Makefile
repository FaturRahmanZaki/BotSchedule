.PHONY: help build run install clean test

help:
	@echo "TurSchedule Bot - Available Commands"
	@echo ""
	@echo "make install    - Download and install dependencies"
	@echo "make build      - Build the application"
	@echo "make run        - Run the bot"
	@echo "make dev        - Run in development mode (with auto-reload)"
	@echo "make clean      - Clean build artifacts"
	@echo "make test       - Run tests"
	@echo ""

install:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed"

build: install
	@echo "🔨 Building application..."
	go build -o turschedule main.go
	@echo "✅ Build complete"

run: build
	@echo "🚀 Starting bot..."
	./turschedule

dev: install
	@echo "🚀 Starting bot in development mode..."
	@echo "Note: This requires 'air' for hot reload"
	go run main.go

clean:
	@echo "🧹 Cleaning up..."
	rm -f turschedule turschedule.exe
	go clean
	@echo "✅ Cleanup complete"

test:
	@echo "🧪 Running tests..."
	go test ./...

fmt:
	@echo "✨ Formatting code..."
	go fmt ./...

lint:
	@echo "🔍 Running linter..."
	@echo "Note: This requires 'golangci-lint' to be installed"
	golangci-lint run ./...
