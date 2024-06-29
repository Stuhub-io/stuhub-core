# --- Tooling & Variables ----------------------------------------------------------------
include ./misc/make/tools.Makefile

# --- ENVS - DATABASE ENVs -----------------------------------------------------------------------
ifneq (,$(wildcard build/local/postgres/.env))
    include build/local/postgres/.env
    export
endif

# ~~~ Development Environment ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
setup:
	@ echo "Setting up the project dependencies ..."
	@ make install-deps
	@ make deps
	@ make down
	@ make up
	@ make migrate-up

up: # Startup / Spinup Docker Compose and air
	@ make dev-env  			

down: docker-stop               ## Stop Docker
destroy: docker-teardown clean  ## Teardown (removes volumes, tmp files, etc...)

install-deps: install-golangci-lint install-air install-golang-migrate install-gorm-gentool

deps:
	@ echo "Required Tools Are Available"

dev-env:
	@ docker-compose -f local.yml up --build -d --remove-orphans

docker-stop:
	@ docker-compose -f local.yml down

docker-teardown:
	@ docker-compose -f local.yml down --remove-orphans -v

# ~~~ Database ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
POSTGRESQL_DSN = postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@127.0.0.1:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
migrate-up:
	@ migrate  -database $(POSTGRESQL_DSN) -path=misc/migrations --verbose up

migrate-down:
	@ migrate  -database $(POSTGRESQL_DSN) -path=misc/migrations --verbose down

migrate-create: 
	@ read -p "Please provide name for the migration: " Name; \
    migrate create -ext sql -dir misc/migrations $${Name}

migrate-drop:
	@ migrate  -database $(POSTGRESQL_DSN) -path=misc/migrations drop

gen-struct:
	@ gentool -c ./gen.yaml


open-db: # CLI for open db using tablePlus only
	@ open $(POSTGRESQL_DSN)

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
	

# ~~~ Testing ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
test:
	go test -v -cover ./...


# ~~~ Swagger ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
swagger:
	@echo Starting swagger generating
	swag init -g **/**/*.go
	make swag-format

swag-format:
	swag fmt

.PHONY: migrate-up migrate-down migrate-create migrate-drop