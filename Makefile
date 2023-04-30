SERVICE := chatgw

all: help

help: ## Show help messages
	@echo "Container - ${SERVICE} "
	@echo
	@echo "Usage:\tmake COMMAND"
	@echo
	@echo "Commands:"
	@sed -n '/##/s/\(.*\):.*##/  \1#/p' ${MAKEFILE_LIST} | grep -v "MAKEFILE_LIST" | column -t -c 2 -s '#'

build: ## Build
	docker compose build --no-cache

up: ## Create and start containers
	docker compose up -d --force-recreate

start: ## Create and start containers
	docker compose start

stop: ## Stop containers
	docker compose stop

logs: ## Check logs
	docker compose logs ${SERVICE}

bash: ## Execute a command in a running container
	docker compose exec ${SERVICE} /bin/bash

.PHONY: build
