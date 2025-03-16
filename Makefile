PG_CLIENT_VERSION=17

.PHONY: all install-dev install-pg-client
PROJECT_DIRECTORY := $(shell pwd)

new-migration:
	./.bin/migrate create -ext .sql -dir migrations $(name)

## Default installation target
all: install-dev

## Install dev tools
install-dev: install-pg-client install-go-migrate install-staticcheck

install-pg-client:
	@echo "Installing PostgreSQL client for Linux..."
	sudo apt update && sudo apt install -y postgresql-client postgresql-client-common

# I would've included it inside go tool but honestly it has too many dependencies
# that I don't want in the go project
install-go-migrate:
	@echo "Installing golang-migrate for Linux..."
	curl -L --silent https://github.com/golang-migrate/migrate/releases/download/v4.18.2/migrate.linux-amd64.tar.gz | tar xvz --directory .bin migrate

install-staticcheck:
	GOBIN="${PROJECT_DIRECTORY}/.bin" go install honnef.co/go/tools/cmd/staticcheck@v0.6.1
	chmod +x "${PROJECT_DIRECTORY}/.bin/staticcheck"

unit-tests:
	GOEXPERIMENT=synctest go test -v ./...

integration-tests:
	chmod +x ./scripts/integration_tests.sh
	./scripts/integration_tests.sh
