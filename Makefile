include .env
export

###

.PHONY: dockerUp dockerDown test

###

_dockerUp:
	docker-compose up -d

_waitForPG:
	$(shell pwd)/scripts/wait_for_postgres.sh

_migrate:
	atlas schema apply --auto-approve -u "postgres://localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?user=$(POSTGRES_USER)&sslmode=disable" --to file://migrate/postgres

dockerUp: _dockerUp _waitForPG _migrate

###

dockerDown:
	docker-compose down

###

TESTS=

_test: dockerUp
	@if [ -z $(TESTS) ]; then \
  		godotenv go test ./... -race; \
  	else \
		godotenv go test $(TESTS) -race; \
	fi

test: _test dockerDown
