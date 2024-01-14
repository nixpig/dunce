ifneq (,$(wildcard .env))
	include .env
	export
endif

APP_PACKAGE_PATH := ./cmd/app
APP_BINARY_NAME := app

.PHONY: tidy
tidy: 
	go fmt ./...
	go mod tidy -v

.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

.PHONY: test
test: export APP_ENV=test
test: 
	go test -v -race -buildvcs ./...

.PHONY: build_app
build_app:
	go build -o tmp/bin/${APP_BINARY_NAME} ${APP_PACKAGE_PATH}

.PHONY: build
build: build_app

.PHONY: dev_app
dev_app: export APP_ENV=development
dev_app: 
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build_app" \
		--build.bin "tmp/bin/${APP_BINARY_NAME}" \
		--build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go" \
		--misc.clean_on_exit "true"

clean:
	rm -rf bin tmp 

.PHONY: dev
dev: export APP_ENV=development
dev:
	make -j1 dev_app

.PHONY: migrate_up
migrate_up:
	migrate -path db/migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${POSTGRES_DB}?sslmode=disable up

.PHONY: migrate_down
migrate_down:
	migrate -path db/migrations -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${POSTGRES_DB}?sslmode=disable down

.PHONY: env
env: 
	# Echos out environment variables
	APP_PACKAGE_PATH=${APP_PACKAGE_PATH}
	APP_BINARY_NAME=${APP_BINARY_NAME}
	APP_ENV=${APP_ENV}
	POSTGRES_USER=${POSTGRES_USER}
	POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
	POSTGRES_DB=${POSTGRES_DB}
	DATABASE_HOST=${DATABASE_HOST}
	DATABASE_PORT=${DATABASE_PORT}

