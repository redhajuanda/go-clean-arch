include .env
export

psql:
	@PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -U ${DB_USERNAME} -d ${DB_NAME}

.PHONY: migrate-up
migrate-up:
	@go run main.go migrate up

.PHONY: migrate-down
migrate-down:
	@go run main.go migrate down

.PHONY: migrate-fresh
migrate-fresh:
	@go run main.go migrate fresh

.PHONY: migrate-new
migrate-new: ## create a new database migration
	
	@read -p "Enter the name of the new migration: " name; \
	CGO_ENABLED="0" go install github.com/rubenv/sql-migrate/...@latest; \
	sql-migrate new $${name}

test:
	go test -v -cover ./...

build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${APP_NAME} -a -installsuffix cgo -ldflags '-w'

.PHONY: build-docker
build-docker: ## build the service as a docker image
	docker build -f Dockerfile -t application .

run:
	@go run main.go api

cron:
	@go run main.go cron

ssh:
	@ssh mcustomer@172.31.143.10 