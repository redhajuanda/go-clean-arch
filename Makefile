include .env
export

psql:
	@PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -U ${DB_USERNAME} -d ${DB_NAME}

.PHONY: migrate-up
migrate-up:
	@go run main.go migrate up

.PHONY: migrate-fresh
migrate-fresh:
	@go run main.go migrate fresh

.PHONY: migrate-new
migrate-new: ## create a new database migration
	@read -p "Enter the name of the new migration: " name; \
	sql-migrate new $${name}

run:
	@go run main.go rest

cron:
	@go run main.go cron

ssh:
	@ssh mcustomer@172.31.143.10 