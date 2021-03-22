DEST_DIR = ./bin
APP_FOLDER = $(HOME)/.kss
CONFIG = $(if \
        $(filter ${MODE},dev),\
                $(APP_FOLDER)/dev_config.yml,\
        $(else)\
                $(APP_FOLDER)/config.yml)

all: build

build: src/main
	go build -o $(DEST_DIR)/server ./src/main

move_config: create_folders
	cp -a extra/* $(APP_FOLDER)

create_folders:
	mkdir -p $(APP_FOLDER)/logs

create_database: DB_HOST = $(call umyml,db_host)
create_database: DB_PORT = $(call umyml,db_port)
create_database: DB_NAME = $(call umyml,db_name)
create_database: NEW_DB_USER = $(call umyml,db_username)
create_database: NEW_USER_PASSWORD = $(call umyml,db_password)
create_database:
	@read -p "Enter admin username for postgres database: " DB_USER; \
	export isDev=`sh -c 'if [ $$MODE == dev ]; then echo 1; else echo 0; fi'`; \
    psql -h $(DB_HOST) -p $(DB_PORT) -U $$DB_USER -d postgres -v db="$(DB_NAME)" -v username="$(NEW_DB_USER)" -v dev=$$isDev -v password="$(NEW_USER_PASSWORD)" -f database-setup.sql

setup: build move_config create_database
	go mod tidy

define umyml
$(shell grep $1 $(CONFIG) | tail -n1 | awk '{print $$2}' | sed -e 's/["]//g')
endef