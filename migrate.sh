#!/bin/bash

# Configuration
CONTAINER_NAME="atlas-dev"
POSTGRES_PASSWORD="pg"
DEV_PORT="5433"

echo "Starting temporary PostgreSQL container..."
docker run --name $CONTAINER_NAME \
	-e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
	-d -p $DEV_PORT:5432 postgres

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to start..."
until docker exec $CONTAINER_NAME pg_isready -h 127.0.0.1; do
	echo "Waiting for PostgreSQL..."
	sleep 1
done

echo "here $GOOSE_DBSTRING"
echo "Running Atlas migration..."
atlas schema apply \
	--to "file://sql/schema.sql" \
	--dev-url "postgres://postgres:$POSTGRES_PASSWORD@127.0.0.1:$DEV_PORT/postgres?sslmode=disable" \
	-u "$GOOSE_DBSTRING"

# Capture the exit status of atlas
ATLAS_EXIT_CODE=$?

echo "Cleaning up: stopping and removing container..."
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME
# Exit with the same status as atlas
exit $ATLAS_EXIT_CODE
