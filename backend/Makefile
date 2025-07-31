#===================#
#== Env Variables ==#
#===================#

.DEFAULT_GOAL := run
APP_DIR = cmd/shortly
DOCKER_COMPOSE_FILE ?= docker-compose.yaml
MIGRATION_NAME ?= new_migration
ROLLBACK_COUNT ?= 1

MAKEFLAGS += --output-sync=target
SHELL := /bin/bash

#==================#
#== Color Values ==#
#==================#

GREEN := \033[0;32m
RED := \033[0;31m
YELLOW := \033[0;33m
RESET := \033[0m

#========================#
#== DATABASE MIGRATION ==#
#========================#

migrate-up:
	@printf "${YELLOW}Running Migrations...${RESET}\n" 
	@docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate up && \
		printf "${GREEN}Migrations ran successfully.${RESET}\n" || \
		printf "${RED}Migrations failed.${RESET}\n"

migrate-down: 
	@printf "${YELLOW}Rolling back migrations...${RESET}\n"
	@docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate down ${ROLLBACK_COUNT} && \
		printf "${GREEN}Rolled back last ${ROLLBACK_COUNT} migrations.${RESET}\n" || \
		printf "${RED}Failed to roll back ${ROLLBACK_COUNT} migrations.${RESET}\n"

migrate-create:
		@printf "${YELLOW}Creating migrations...${RESET}\n"
	@docker compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate create -ext sql -dir /migrations -seq ${MIGRATION_NAME} && \
		printf "${GREEN}Migrations created for ${MIGRATION_NAME}.${RESET}\n" || \
		printf "${RED}Failed to create migrations.${RESET}\n"

shell-db:
	@docker compose -f ${DOCKER_COMPOSE_FILE} exec postgres psql -U postgres -d postgres

shell-redis:
	@docker compose -f ${DOCKER_COMPOSE_FILE} exec redis redis-cli

start-services:
	@printf "${YELLOW}Starting docker services...${RESET}\n"
	@docker compose -f ${DOCKER_COMPOSE_FILE} up -d && \
		printf "${GREEN}Services started.${RESET}\n" || \
		printf "${RED}Failed to start services.${RESET}\n"

stop-services:
	@printf "${YELLOW}Stopping docker services...${RESET}\n"
	@docker compose -f ${DOCKER_COMPOSE_FILE} down --remove-orphans && \
		printf "${GREEN}Services stopped.${RESET}\n" || \
		printf "${RED}Failed to stop services.${RESET}\n"

build:
	@printf "${YELLOW}Building go app...${RESET}\n"
	@go build -C ${APP_DIR} -o ../../bin/main && \
		printf "${GREEN}Golang binary built.${RESET}\n" || \
		printf "${RED}Failed to build binary.${RESET}\n"

run: start-services migrate-up build 
	@printf "${YELLOW}Running app...${RESET}\n"
	@./bin/main

dev: start-services migrate-up
	@go run github.com/air-verse/air -c .air.toml

dev-restart:
	@printf "${YELLOW}Quick restart with existing containers...${RESET}\n"
	@go run github.com/air-verse/air -c .air.toml

swagger:
	@printf "${YELLOW}Generating swagger documentation...${RESET}\n"
	@which swag >/dev/null || go install github.com/swaggo/swag/cmd/swag@latest
	@swag init -g cmd/shortly/main.go -o ./api && \
		printf "${GREEN}Swagger documentation generated.${RESET}\n" || \
		printf "${RED}Failed to generate swagger documentation.${RESET}\n"
