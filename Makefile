ifneq ($(wildcard .env),)
	include .env
endif

export $(shell sed 's/=.*//' .env)

migration_path ?= ./cmd/migration

create-migrate:
	@migrate create -ext sql -dir $(migration_path) $(sql)


force-migrate:
	@migrate -path $(migration_path) -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=$(POSTGRES_SSLMODE)" force 1


up-migrate:
	@migrate -path $(migration_path) -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=$(POSTGRES_SSLMODE)" -verbose up


down-migrate:
	@migrate -path $(migration_path) -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=$(POSTGRES_SSLMODE)" -verbose down