SHELL := /bin/bash
DEST_DIR = ./bin
APP_FOLDER = $(HOME)/.kss

ifneq ("$(wildcard ./.env)","")
    CONFIG := ./.env
else
	CONFIG := ./.example_env
endif

include $(CONFIG)

.DEFAULT_GOAL := build

define shutdown_server
curl http://${ADDR}:${PORT}/dev/shutdown-server
endef

build: src/main
	go build -o $(DEST_DIR)/server ./src/main

create_folders:
	mkdir -p $(APP_FOLDER)/logs
	$(info >>>> Created folders)

create_database:
	bash ./initdb.sh local

tidy:
	go mod tidy

setup: tidy build create_folders
	$(info >>>> Setup complete!)

setup_docker:
	sudo docker-compose build
	sudo docker-compose up

setup_bin: setup
	cp ./example_env ./bin/.env


tests: tidy
	$(info $(CONFIG))
	go test -count=1 \
		$(shell find -name "*_test.go" -exec dirname {} \; | uniq) || $(call shutdown_server)

	$(call shutdown_server)
