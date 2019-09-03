#!/usr/bin/env bash

cd ..
export SAVINGS_DB_HOST=localhost
export SAVINGS_DB_NAME=savings-test
export SAVINGS_DB_USER=afruizc
go test ./...

