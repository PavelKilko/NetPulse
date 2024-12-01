#!/bin/sh

echo "Loading environment variables from .env..."
export $(grep -v '^#' /root/.env | xargs)

# Debug print to ensure environment variables are set correctly
echo "POSTGRES_USER=$POSTGRES_USER"
echo "POSTGRES_PASSWORD=$POSTGRES_PASSWORD"
echo "POSTGRES_DB=$POSTGRES_DB"

echo "Waiting for PostgreSQL..."
/root/scripts/wait-for-it.sh postgres 5432 -- echo "PostgreSQL is up"

echo "Waiting for Redis..."
/root/scripts/wait-for-it.sh redis 6379 -- echo "Redis is up"

echo "Checking if PostgreSQL is ready to accept connections..."
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "postgres" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q' >/dev/null 2>&1; do
  echo "PostgreSQL is not ready yet..."
  sleep 2
done

echo "PostgreSQL is ready to accept connections."

echo "Running migrations..."
migrate -path ./db/migration -database "$DATABASE_URL" up

echo "Starting NetPulse..."
./netpulse-backend
