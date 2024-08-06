#!/bin/bash
set -e

if [ -z "$POSTGRES_USER" ]; then
  echo "Error: POSTGRES_USER is not set."
  exit 1
fi

if [ -z "$POSTGRES_DB" ]; then
  echo "Error: POSTGRES_DB is not set."
  exit 1
fi

echo "POSTGRES_USER is set to ${POSTGRES_USER}"
echo "POSTGRES_DB is set to ${POSTGRES_DB}"

 
echo "Creating database ${POSTGRES_DB}..."
psql -U "$POSTGRES_USER" -c "CREATE DATABASE ${POSTGRES_DB};"