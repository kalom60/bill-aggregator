#!/bin/bash
set -e

# Wait for PostgreSQL to be ready (as postgres user)
until PGPASSWORD=${POSTGRES_PASSWORD} psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -c '\q'; do
  >&2 echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

# Execute SQL as postgres user
PGPASSWORD=${POSTGRES_PASSWORD} psql -v ON_ERROR_STOP=1 -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" <<-EOSQL
    -- Create roles with passwords from environment
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'user_service') THEN
            CREATE ROLE user_service LOGIN PASSWORD '${USER_SERVICE_PASSWORD}';
        END IF;
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'provider_service') THEN
            CREATE ROLE provider_service LOGIN PASSWORD '${PROVIDER_SERVICE_PASSWORD}';
        END IF;
        IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'account_service') THEN
            CREATE ROLE account_service LOGIN PASSWORD '${ACCOUNT_SERVICE_PASSWORD}';
        END IF;
    END
    \$\$;

    -- Create schemas
    CREATE SCHEMA IF NOT EXISTS user_service;
    CREATE SCHEMA IF NOT EXISTS provider_service;
    CREATE SCHEMA IF NOT EXISTS account_service;

    -- Grant permissions
    GRANT ALL PRIVILEGES ON SCHEMA user_service TO user_service;
    GRANT ALL PRIVILEGES ON SCHEMA provider_service TO provider_service;
    GRANT ALL PRIVILEGES ON SCHEMA account_service TO account_service;

    -- Grant permissions on all tables in the user_service schema
    GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA user_service TO user_service;
    -- Grant permissions on all tables in the provider_service schema
    GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA provider_service TO provider_service;
    -- Grant permissions on all tables in the account_service schema
    GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA account_service TO account_service;

    -- Grant permissions on future tables in the user_service schema
    ALTER DEFAULT PRIVILEGES IN SCHEMA user_service GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO user_service;
    -- Grant permissions on future tables in the provider_service schema
    ALTER DEFAULT PRIVILEGES IN SCHEMA provider_service GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO provider_service;
    -- Grant permissions on future tables in the accounut_service schema
    ALTER DEFAULT PRIVILEGES IN SCHEMA account_service GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO account_service;
EOSQL

>&2 echo "Database initialization complete!"

