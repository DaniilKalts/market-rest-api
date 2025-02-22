.PHONY: build run clean

build:
	go build -o market-rest-api

run: build
	./market-rest-api

clean:
	rm -f market-rest-api
