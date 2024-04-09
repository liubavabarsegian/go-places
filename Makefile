NETWORK_NAME := go-places-network

start:
	docker-compose build
	docker-compose up -d
	make connect-network

stop:
	docker-compose down

restart:
	docker-compose down
	make start

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
