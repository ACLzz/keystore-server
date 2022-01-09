#!/bin/bash
if [ -f ./.env ]; then
  export CONFIG=./.env
else
  if [[ $1 == "local" ]]; then
    export CONFIG=./.example_env
  else
    export CONFIG=./.docker_env
  fi
fi

echo ">>>> Reading config from $CONFIG"
source $CONFIG
export $(cut -d= -f1 $CONFIG)
export isDev=`bash -c 'if [[ $$MODE == test ]]; then echo 1; else echo 0; fi'`
if [[ $1 == "local" ]]; then
  psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U postgres -d postgres -v db=$POSTGRES_NAME -v username=$POSTGRES_USERNAME -v dev=$isDev -v password=$POSTGRES_PASSWORD -f database-setup.sql && \
    echo ">>>> Database setup done!"
else
  psql -U postgres -v db=$POSTGRES_NAME -v username=$POSTGRES_USERNAME -v dev=$isDev -v password=$POSTGRES_PASSWORD -f database-setup.sql && \
    echo ">>>> Database setup done!"
fi
