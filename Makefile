.PHONY: create-db build run

# If you have a different database name, replace 'market' with your database name
create-db:
	psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'market'" | grep -q 1 || \
	psql -U postgres -c "CREATE DATABASE market"

build: create-db
	go build -o market-rest-api

run: build
	./market-rest-api
