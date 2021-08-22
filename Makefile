SHELL := /bin/bash
DEST_DIR = ./bin
APP_FOLDER = $(HOME)/.kss

ifneq ("$(wildcard $(PATH_TO_FILE))","")
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

create_database:
	source $(CONFIG); \
	export $(cut -d= -f1 $(CONFIG)); \
	if [ $$MODE != test ]; \
	then \
		read -p "Enter admin username for postgres database: " DB_USER; \
	else \
	  export DB_USER=postgres; \
	fi; \
	export isDev=`sh -c 'if [ $$MODE = test ]; then echo 1; else echo 0; fi'`; \
    psql -h $$DB_HOST -p $$DB_PORT -U $$DB_USER -d postgres -v db=$$DB_NAME -v username=$$DB_USERNAME -v dev=$$isDev -v password=$$DB_PASSWORD -f database-setup.sql

tidy:
	go mod tidy

setup: tidy build create_folders create_database
	go mod tidy

tests:
	$(info $(CONFIG))
	cd bin;\
	./server &
	go test -count=1 \
		$(shell find -name "*_test.go" -exec dirname {} \; | uniq) || $(call shutdown_server)

	$(call shutdown_server)
