#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER diceroller;
    CREATE DATABASE diceroller;
    GRANT ALL PRIVILEGES ON DATABASE diceroller TO diceroller;
EOSQL
