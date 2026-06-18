.PHONY: build run server cli tidy test setup-labels

build: server cli

server:
	go build -o bin/feedback-server ./cmd/server

cli:
	go build -o bin/rcli ./cmd/rcli

run: server
	@echo "Starting feedback server..."
	@set -a && [ -f .env ] && . ./.env; set +a && ./bin/feedback-server

tidy:
	go mod tidy

test:
	go test ./...

setup-labels:
	@set -a && [ -f .env ] && . ./.env; set +a && bash scripts/setup-github-labels.sh

install-cli:
	go install ./cmd/rcli

deploy:
	bash scripts/deploy-gcp.sh

e2e:
	bash scripts/e2e-test.sh