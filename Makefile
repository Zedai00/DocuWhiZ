.PHONY: all build run clean frontend backend

# Default port for the server
PORT ?= 8000

# Require the GEMINI_API_KEY environment variable to be set
ifeq ($(GEMINI_API_KEY),)
$(error GEMINI_API_KEY is not set. Please set it via "export GEMINI_API_KEY=your_key_here" or by providing it on the command line: "make run GEMINI_API_KEY=your_key_here")
endif

# Default target runs the application
all: run

# Build the frontend and backend
build: frontend backend

# Build the frontend using npm
frontend:
	@echo "--- Building frontend ---"
	cd frontend && npm install && npm run build

# Build the Go backend
backend:
	@echo "--- Building backend ---"
	# The key change is here: changing the directory before running Go commands
	cd backend && go mod tidy && go build -o server .
	cp -r frontend/dist backend/

# Run the Go server with the required environment variables
run: build
	@echo "--- Starting server on port $(PORT) ---"
	cd backend && GEMINI_API_KEY=$(GEMINI_API_KEY) PORT=$(PORT) ./server

# Clean up built files
clean:
	@echo "--- Cleaning up built files ---"
	rm -f backend/server
	rm -rf backend/dist
	rm -rf frontend/node_modules
	rm -rf frontend/dist

