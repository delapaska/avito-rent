#!/bin/bash
set -e

host="$DB_HOST"
port="$DATABASE_PORT"
user="$DB_USER"
dbname="$DB_NAME"

until pg_isready -h "$host" -p "$port" -U "$user" -d "$dbname"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec "$@"
