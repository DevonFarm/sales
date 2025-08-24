.DEFAULT_GOAL:=help
SHELL:=/usr/bin/env bash

help: ## Display this help
	@echo "Usage:"
	@echo -e "  make \x1b[1;36m<target>\x1b[0m"
	@echo
	@echo "Targets:"
	@sed \
		-e '/^[a-zA-Z0-9_\-]*:.*##/!d'                \
		-e 's/:.*##\s*/:/'                            \
		$(MAKEFILE_LIST)                            | \
		column -c2 -t -s :                          | \
		sed 's/^\([^ ]*\)/  \x1b[1;36m\1\x1b[0m/'

.PHONY: db-migrate
db-migrate: ## Runs database migrations using golang-migrate
	@echo "Running db migrations..."
  @export COCKROACH_DSN="$(grep COCKROACH_DSN .env | cut -d '=' -f2- | tr -d '"')"
  migrate -database "${COCKROACH_DSN}" -path ./db/migrations up
