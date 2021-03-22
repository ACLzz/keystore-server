DEST_DIR = ./bin
APP_FOLDER = $(HOME)/.kss
CONFIG = $(if \
        $(filter ${MODE},dev),\
                $(APP_FOLDER)/dev_config.yml,\
		$(if \
		$(filter ${MODE},test),\
				$(APP_FOLDER)/test_config.yml,\
        $(else)\
                $(APP_FOLDER)/config.yml))

.DEFAULT_GOAL := build

define umyml
$(shell grep $1 $(CONFIG) | head -n1 | awk '{print $$2}' | sed -e 's/["]//g')
endef

define shutdown_server
curl http://$(call umyml,addr):$(call umyml,port)/dev/shutdown-server
endef

build: src/main
	go build -o $(DEST_DIR)/server ./src/main

move_config: create_folders
	$(info $(CONFIG))
	cp -a extra/* $(APP_FOLDER)

create_folders:
	mkdir -p $(APP_FOLDER)/logs

create_database: DB_HOST = $(call umyml,db_host)
create_database: DB_PORT = $(call umyml,db_port)
create_database: DB_NAME = $(call umyml,db_name)
create_database: NEW_DB_USER = $(call umyml,db_username)
create_database: NEW_USER_PASSWORD = $(call umyml,db_password)
create_database:
	echo ; \
	if [ $$MODE != test ]; \
	then \
		read -p "Enter admin username for postgres database: " DB_USER; \
	else \
	  export DB_USER=postgres; \
	fi; \
	export isDev=`sh -c 'if [ $$MODE = test ]; then echo 1; else echo 0; fi'`; \
    psql -h $(DB_HOST) -p $(DB_PORT) -U $$DB_USER -d postgres -v db="$(DB_NAME)" -v username="$(NEW_DB_USER)" -v dev=$$isDev -v password="$(NEW_USER_PASSWORD)" -f database-setup.sql

setup: build move_config create_database
	go mod tidy

tests:
	$(info $(CONFIG))
	cd bin;\
	./server &
	go test -count=1 \
		$(shell find -name "*_test.go" -exec dirname {} \; | uniq) || $(call shutdown_server)

	$(call shutdown_server)