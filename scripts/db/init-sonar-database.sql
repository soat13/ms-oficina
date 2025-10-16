-- Initialization script for SonarQube database
-- Executed only once when the container is created

SELECT 'CREATE DATABASE sonarqube'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'sonarqube')\gexec

