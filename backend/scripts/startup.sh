#!/bin/sh

set -e

echo "Loading environment variables from .env..."
export $(grep -v '^#' /root/.env | xargs)

# Debug print to ensure environment variables are set correctly
echo "PORT=$PORT"
echo "DATABASE_URL=$DATABASE_URL"
echo "REDIS_URL=$REDIS_URL"
echo "RABBITMQ_URL=$RABBITMQ_URL"
echo "MONGODB_URL=$MONGODB_URL"
echo "POSTGRES_USER=$POSTGRES_USER"
echo "POSTGRES_PASSWORD=$POSTGRES_PASSWORD"
echo "POSTGRES_DB=$POSTGRES_DB"

# Function to wait for a service to be available
wait_for_service() {
  host=$1
  port=$2
  echo "Waiting for $host:$port..."
  until nc -z "$host" "$port"; do
    sleep 2
  done
  echo "$host:$port is up"
}

# Wait for PostgreSQL
wait_for_service postgres 5432

# Ensure PostgreSQL is ready to accept connections
echo "Checking if PostgreSQL is ready to accept connections..."
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "postgres" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q' >/dev/null 2>&1; do
  echo "PostgreSQL is not ready yet..."
  sleep 2
done
echo "PostgreSQL is ready to accept connections."

# Wait for Redis
wait_for_service redis 6379

# Wait for RabbitMQ
wait_for_service rabbitmq 5672

# Wait for MongoDB
wait_for_service mongodb 27017

# Run migrations
echo "Running migrations..."
migrate -path ./db/migration -database "$DATABASE_URL" up

# Start the application
echo "Starting NetPulse..."
exec ./netpulse-backend
