.PHONY: create-db build run docker-clean docker-run

# Without Docker

# If you have a different database name, replace 'market' with your database name
create-db:
	psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'market'" | grep -q 1 || \
	psql -U postgres -c "CREATE DATABASE market"

build: create-db
	go build -o market-rest-api ./cmd/market-rest-api

run: build
	./market-rest-api

# With Docker

docker-clean:
	docker-compose down --remove-orphans
	docker system prune -f

docker-run: docker-clean
	docker-compose up --build
