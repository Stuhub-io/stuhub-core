# Load environment variables from .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

POSTGRESQL_USER ?= postgres
POSTGRESQL_PASSWORD ?= password
POSTGRESQL_ADDRESS ?= 127.0.0.1:5432
POSTGRESQL_DATABASE ?= stuhub
POSTGRESQL_CONTAINER_NAME ?= postgres-db

# ~~~ Development Environment ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
up: dev-env dev-air   			## Startup / Spinup Docker Compose and air
down: docker-stop               ## Stop Docker
destroy: docker-teardown clean  ## Teardown (removes volumes, tmp files, etc...)

dev-env:
	@ docker-compose up -d --build pgsql

dev-air: 
	@ air

docker-stop:
	@ docker-compose down

docker-teardown:
	@ docker-compose down --remove-orphans -v


# ~~~ Modules support ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
tidy:
	go mod tidy

deps-reset:
	git checkout -- go.mod
	go mod tidy

deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy

deps-cleancache:
	go clean -modcache


# ~~~ Code Actions ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
lint:
	@echo Starting linters
	golangci-lint run ./...


# ~~~ Database ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
createdb:
	docker exec -it =$(POSTGRESQL_CONTAINER_NAME) createdb --username=$(POSTGRESQL_USER) --owner=$(POSTGRESQL_USER) $(POSTGRESQL_DATABASE)

dropdb:
	docker exec -it =$(POSTGRESQL_CONTAINER_NAME) dropdb --username==$(POSTGRESQL_USER) $(POSTGRESQL_DATABASE)

POSTGRESQL_DSN := "postgresql://$(POSTGRESQL_USER):$(POSTGRESQL_PASSWORD)@$(POSTGRESQL_ADDRESS)/$(POSTGRESQL_DATABASE)?sslmode=disable"

migrate-up:
	migrate  -database $(POSTGRESQL_DSN) -path=misc/migrations --verbose up

migrate-down:
	migrate  -database $(POSTGRESQL_DSN) -path=misc/migrations --verbose down

migrate-create: 
	@ read -p "Please provide name for the migration: " Name; \
    migrate create -ext sql -dir misc/migrations $${Name}

migrate-drop:
	migrate  -database $(POSTGRESQL_DSN) -path=misc/migrations drop
	

# ~~~ Testing ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
test:
	go test -v -cover ./...


# ~~~ Swagger ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
swagger:
	@echo Starting swagger generating
	swag init -g **/**/*.go

.PHONY: migrate-up migrate-down migrate-create migrate-drop