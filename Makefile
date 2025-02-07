PG_CLIENT_VERSION=17

.PHONY: all install-dev install-pg-client 

## Default installation target
all: install-dev

## Install dev tools
install-dev: install-pg-client

install-pg-client:
	@echo "Installing PostgreSQL client for Linux..."
	sudo apt update && sudo apt install -y postgresql-client postgresql-client-common
