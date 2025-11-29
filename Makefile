run:
	go run cmd/api/main.go

generate-sql:
	sqlc generate -f sqlc.yaml

.PHONY: migration

migration:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir db/migrations -seq "$$(date +%Y%m%d_%H%M%S)_$${name}"

swagger:
	swag init -g cmd/api/main.go

include .env
export $(shell sed 's/=.*//' .env)
.PHONY: migrate-up

migrate-down:
	@read -p "version: " version; \
		migrate -path db/migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		down $${version}

migrate-force:
		@read -p "version: " version; \
		migrate -path db/migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		force $${version}

migrate-up:
	migrate -path db/migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		up

build:
	swag init -g cmd/main.go -o docs
	go run cmd/main.go
