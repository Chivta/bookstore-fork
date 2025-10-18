#!/bin/bash
set -e

# Create separate databases for each microservice
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Create databases for each service
    CREATE DATABASE bookstore_books;
    CREATE DATABASE bookstore_users;
    CREATE DATABASE bookstore_logs;

    -- Grant privileges
    GRANT ALL PRIVILEGES ON DATABASE bookstore_books TO bookstore;
    GRANT ALL PRIVILEGES ON DATABASE bookstore_users TO bookstore;
    GRANT ALL PRIVILEGES ON DATABASE bookstore_logs TO bookstore;

    \c bookstore_books;
    -- Enable UUID extension for books database
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    \c bookstore_users;
    -- Enable UUID extension for users database
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

    \c bookstore_logs;
    -- Enable UUID extension for logs database
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
EOSQL

echo "Databases created successfully"
