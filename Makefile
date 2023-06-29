# Makefile for leeta_backend

# Directories
CMD_DIR := ./cmd
WORKDIR := $(shell pwd)

# Docker commands
DOCKER := docker
DOCKER_COMPOSE := docker-compose

# MongoDB container
MONGODB_CONTAINER := leeta-mongo1
MONGODB_IMAGE := mongo
MONGODB_PORT := 27017

# Database and user credentials
DB_NAME := leetadb
DB_USER := leeta
DB_PASSWORD := leet
DB_AUTH_MECHANISM := SCRAM-SHA-256

# URLs
SWAGGER_URL := http://localhost:3000/leeta/swagger/index.html

# Targets
.PHONY: all start stop run_app stop_app

all: start

start: check_docker check_mongodb create_user check_database get_dependencies init_swagger run_app wait_before_open_browser open_browser
	@echo "To start the application, run 'make run_app'"

stop-mongo:
	@echo "Stopping MongoDB container..."
	@$(DOCKER) stop $(MONGODB_CONTAINER)

run_app:
	@echo "Running the application..."
	@cd $(CMD_DIR) && go run main.go &

stop_app:
	@echo "Stopping the application..."
	@pkill -f "main" || true


#stop_app:
#	@echo "Stopping the application..."
#	@pkill -INT -f "go run main.go" || true

check_docker:
	@echo "Checking if Docker $(DOCKER) is installed..."
	@if ! command -v $(DOCKER); then \
		echo "Docker not found. Please install Docker."; \
		exit 1; \
	fi

check_mongodb:
	@echo "Checking if MongoDB is installed and running..."
	@if ! $(DOCKER) ps -a --format '{{.Names}}' | grep -q $(MONGODB_CONTAINER); then \
		echo "MongoDB container not found. Installing MongoDB on Docker..."; \
		$(DOCKER) run -d -p $(MONGODB_PORT):$(MONGODB_PORT) --name $(MONGODB_CONTAINER) --env MONGO_INITDB_DATABASE=$(DB_NAME) $(MONGODB_IMAGE); \
	elif ! $(DOCKER) ps -f "name=$(MONGODB_CONTAINER)" --format '{{.Names}}' | grep -q $(MONGODB_CONTAINER); then \
		echo "MongoDB container found but not running. Starting MongoDB container..."; \
		$(DOCKER) start $(MONGODB_CONTAINER); \
	fi

create_user:
	@echo "Creating user $(DB_USER) for admin database if not exists..."
	@if $(DOCKER) exec $(MONGODB_CONTAINER) mongo admin --quiet --eval 'db.getUsers().forEach(function(user) { if (user.user == "$(DB_USER)") { quit(0); } }); quit(1);'; then \
		echo "User $(DB_USER) does not exist in the admin database. Creating user..."; \
		$(DOCKER) exec -it $(MONGODB_CONTAINER) mongo admin --eval 'db.createUser({ user: "$(DB_USER)", pwd: "$(DB_PASSWORD)", roles: ["root"] })'; \
	else \
		echo "User $(DB_USER) already exists in the admin database. Skipping user creation."; \
	fi

check_database:
	@echo "Checking if database $(DB_NAME) exists..."
	@if ! $(DOCKER) exec $(MONGODB_CONTAINER) mongosh --quiet --eval 'db.getMongo().getDBNames().includes("$(DB_NAME)")' | grep -q true; then \
		echo "Creating database $(DB_NAME)..."; \
		$(DOCKER) exec $(MONGODB_CONTAINER) mongosh --eval 'use $(DB_NAME); db.runCommand({ ping: 1 })'; \
	else \
		echo "Database $(DB_NAME) already exists."; \
	fi


get_dependencies:
	cd $(WORKDIR)/cmd/
	go install github.com/swaggo/swag/cmd/swag@latest

init_swagger:
	@echo "Running Swagger initialization..."
	@echo "Fetching swagger dependency $(GOPATH)"
	go generate ./...

open_browser:
	@echo "Opening browser..."
	@open $(SWAGGER_URL)

wait_before_open_browser:
	@echo "Waiting for 2 seconds before opening the browser..."
	@sleep 2

generate_keys:
	openssl genrsa -out private.key 2048
	openssl rsa -in private.key -pubout -out public.key
	@PRIVATE_KEY="$$(cat private.key)"; \
	PUBLIC_KEY="$$(cat public.key)"; \
	echo "PRIVATE_KEY=\"$$PRIVATE_KEY\"" > local.env; \
	echo >> local.env; \
    echo >> local.env; \
	echo "PUBLIC_KEY=\"$$PUBLIC_KEY\"" >> local.env; \
	rm private.key public.key
