.PHONY: help test test-backend test-frontend coverage lint run-backend run-frontend dev build docker-build docker-run clean

help: ## Show the available commands
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "} {printf "  %-15s %s\n", $$1, $$2}'

test: test-backend test-frontend ## Run all unit tests

test-backend: ## Run Go tests
	cd backend && go test ./...

test-frontend: ## Run frontend tests
	cd frontend && npm test -- --run

coverage: ## Write coverage reports to coverage/
	mkdir -p coverage
	cd backend && go test ./... -coverprofile=../coverage/backend.out && go tool cover -func=../coverage/backend.out | tail -1 && go tool cover -html=../coverage/backend.out -o ../coverage/backend.html
	cd frontend && npx vitest run --coverage --coverage.reportsDirectory=../coverage/frontend
	@echo "Open coverage/backend.html and coverage/frontend/index.html"

lint: ## Check formatting and run linters
	@test -z "$$(gofmt -l backend)" || (echo "gofmt needed:" && gofmt -l backend && exit 1)
	cd backend && go vet ./...
	cd frontend && npx tsc -b && npm run lint

run-backend: ## Run the API on :8080
	cd backend && go run ./cmd/server

run-frontend: ## Run the dev server on :5173 (proxies /api to :8080)
	cd frontend && npm run dev

dev: ## Run backend and frontend together
	$(MAKE) -j2 run-backend run-frontend

build: ## Build the frontend and the backend binary into bin/
	cd frontend && npm run build
	cd backend && go build -o ../bin/server ./cmd/server

docker-build: ## Build the Docker image
	docker build -t sezzle-calculator .

docker-run: docker-build ## Build and run the app on http://localhost:8080
	docker run --rm -p 8080:8080 sezzle-calculator

clean: ## Remove build and coverage output
	rm -rf bin coverage frontend/dist frontend/coverage
