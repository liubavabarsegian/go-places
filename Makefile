NETWORK_NAME := go-places-network

start:
	docker-compose build
	docker-compose up -d

stop:
	docker-compose down

restart:
	docker-compose down
	make start

restart_app:
	docker-compose down go-places-app
	docker-compose build go-places-app
	docker-compose up -d go-places-app

create-network:
	docker-compose down
	@if ! docker network ls | grep -q $(NETWORK_NAME); then \
		echo "Creating Docker network: $(NETWORK_NAME)"; \
		docker network create $(NETWORK_NAME); \
	else \
		echo "Docker network $(NETWORK_NAME) already exists."; \
	fi

connect-network:
	docker network connect go-places-network go-places-app
	docker network connect go-places-network go-places-app-elasticsearch-1

mapping:
	curl -X DELETE "localhost:9200/places"
	curl -X PUT "localhost:9200/places/_mapping" -H "Content-Type: application/json" -d @config/schema.json
