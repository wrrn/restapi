#!/usr/bin/env bash
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
export PGPASSWORD=insecure
PGUSER=tenable
DATABASE=restapi
sudo apt-get update
sudo apt-get install -y postgresql
sudo -u postgres psql -c "DROP DATABASE $DATABASE" || 0
sudo -u postgres psql -c "DROP USER $PGUSER" || 0
sudo -u postgres psql -c "CREATE USER $PGUSER WITH password '$PGPASSWORD'"
sudo -u postgres psql -c "CREATE DATABASE $DATABASE"
psql -U $PGUSER -d $DATABASE -a -f ${DIR}/create_db.sql -h localhost
sudo -u postgres psql -c "DROP DATABASE apitest" || 0
sudo -u postgres psql -c "CREATE DATABASE apitest WITH TEMPLATE $DATABASE"
