
# --- Tooling & Variables ----------------------------------------------------------------
include ./misc/make/tools.Makefile

ENV ?= local # local | production | staggning


include ./build/$(ENV)/postgres/.env
export

# ~~~ Dev without Docker ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
dev:
	@air -c .air.toml

run-cmd:
	@ read -p "Please provide cmd file name: " Name; \
    go run cmd/$${Name}/main.go

# ~~~ Development Environment ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
setup:
	@ echo "Setting up the project dependencies ..."
	@ make install-deps
	@ make deps
	@ make down
	@ make up
	@ make migrate-up

setup-hosted-pg:
	@ echo "Setting up the project dependencies ..."
	@ make install-deps
	@ make deps
	@ make down
	@ make hosted-pg-up
	@ make migrate-prod-up	

up: # Startup / Spinup Docker Compose and air
	@ docker compose -f local.yml up --build -d --remove-orphans

down: docker-stop               ## Stop Docker

destroy: docker-teardown clean  ## Teardown (removes volumes, tmp files, etc...)

install-deps: install-golangci-lint install-air install-golang-migrate install-gorm-gentool

deps:
	@ echo "Required Tools Are Available"

docker-stop:
	@ docker compose -f local.yml down

docker-teardown:
	@ docker compose -f local.yml down --remove-orphans -v

hosted-pg-up: 
	@ docker compose -f local-hosted-pg.yml up --build -d --remove-orphans

# ~~~ Database ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
POSTGRESQL_DSN = postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)

# NOTE: run command with ENV = production for production database
migrate-up:
	@ migrate -database $(POSTGRESQL_DSN) -path=misc/migrations --verbose up

migrate-down:
	@ migrate -database $(POSTGRESQL_DSN) -path=misc/migrations --verbose down

migrate-create:
	@ read -p "Please provide name for the migration: " Name; \
    migrate create -ext sql -dir misc/migrations $${Name}

migrate-force:
	@ migrate -database $(POSTGRESQL_DSN) -path=misc/migrations --verbose force $(VERSION)

migrate-drop:
	@ migrate  -database $(POSTGRESQL_DSN) -path=misc/migrations drop

gen-struct:
	@ gentool -c ./gen.yaml

gen-struct-staging:
	@ gentool -c ./gen-hosted-pg.yaml

open-db: # CLI for open db using tablePlus only
	open $(POSTGRESQL_DSN)

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

.PHONY: migrate-up migrate-down migrate-create migrate-drop migrate-force gen-struct gen-struct-hosted open-db

RELEASE_BRANCH ?= main
BETA_BRANCH ?= develop
DEVELOP_BRANCH ?= develop

.PHONY: release
release: sync-release
	git checkout $(BETA_BRANCH) && git pull origin $(BETA_BRANCH) && \
		git checkout $(RELEASE_BRANCH) && git pull origin $(RELEASE_BRANCH) && \
		git merge $(BETA_BRANCH) --no-edit --no-ff && \
		git push origin $(RELEASE_BRANCH) && \
		git checkout $(DEVELOP_BRANCH) && git push origin $(DEVELOP_BRANCH)

.PHONY: sync-release
sync-release:
	git checkout $(RELEASE_BRANCH) && git pull origin $(RELEASE_BRANCH) && \
		git checkout $(BETA_BRANCH) && git pull origin $(BETA_BRANCH) && \
		git merge $(RELEASE_BRANCH) --no-edit --no-ff