.PHONY: env db_up db_down db_version lint tests pkg_scan build tests_nocov run

DB_URL="postgresql://user:password@host:port/database?sslmode=disable"
APP_NAME="main"

env:
	@if [ ! -f .env.dev ]; then cp .env.example .env.dev; else echo ".env.dev already exists, skipping."; fi
	@if [ ! -f .env.db ]; then cp .env.db.example .env.db; else echo ".env.db already exists, skipping."; fi

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o build/$(APP_NAME) main.go

run:
	@go run main.go

lint:
	@golangci-lint run ./...

tests:
	@go test ./... -coverprofile=artifact/coverage/coverage.out


tests_nocov:
	@go test ./... && rm -rf internal/configs/bootstrap/resource


pkg_scan:
	@osv-scanner --lockfile package-lock.json
	@osv-scanner --lockfile go.mod

db_up:
	@migrate -path ./db/migrations -database "$(DB_URL)" -verbose up

db_down:
	@migrate -path ./db/migrations -database "$(DB_URL)" -verbose down

db_version:
	@migrate -path ./db/migrations -database "$(DB_URL)" -verbose version

