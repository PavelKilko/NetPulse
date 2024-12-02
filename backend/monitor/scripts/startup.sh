#!/bin/sh

set -e

echo "Loading environment variables from .env..."
export $(grep -v '^#' .env | xargs)

# Debug print to ensure environment variables are set correctly
echo "MONGODB_URL=$MONGODB_URL"
echo "RABBITMQ_URL=$RABBITMQ_URL"

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

# Wait for MongoDB
wait_for_service mongodb 27017

# Wait for RabbitMQ
wait_for_service rabbitmq 5672

# Start the Monitor service
echo "Starting Monitor Service..."
exec ./monitor-service
