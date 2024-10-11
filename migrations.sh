#!/usr/bin/bash

# This scipt will run the migrations for the database

if [ ! -d ./migrations ]; then
    echo "Migrations folder not found"
    exit 1
fi

# if param is not passed, then default to up
if [ -z "$1" ]; then
    echo "Running up migrations"
    migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/gooban?sslmode=disable" up
    exit 0
fi

if [ "$1" == "up" ]; then
    echo "Running up migrations"
    migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/gooban?sslmode=disable" up
    exit 0
fi

if [ "$1" == "down" ]; then
    echo "Running down migrations"
    migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/gooban?sslmode=disable" down
    exit 0
fi
