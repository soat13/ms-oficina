#!/bin/bash
set -e

# Initialization script for SonarQube database
# Executed only once when the container is created

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Check if database exists before creating
    SELECT 'CREATE DATABASE sonarqube'
    WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'sonarqube')\gexec
EOSQL

echo "SonarQube database initialization completed successfully"

