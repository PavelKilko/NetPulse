#!/bin/sh
# wait-for-it.sh - Script to wait for a specific port to be available

set -e

host="$1"
shift
port="$1"
shift

until nc -z "$host" "$port"; do
  echo "Waiting for $host:$port..."
  sleep 2
done

# Additional check for PostgreSQL to ensure it's ready to accept connections
if [ "$host" = "postgres" ]; then
  echo "Checking if PostgreSQL is ready to accept connections..."
  until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U $POSTGRES_USER -d $POSTGRES_DB -c '\q' >/dev/null 2>&1; do
    echo "PostgreSQL is not ready yet..."
    sleep 2
  done
  echo "PostgreSQL is ready to accept connections."
fi

exec "$@"
