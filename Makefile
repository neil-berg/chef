up: 
	docker-compose build 
	docker-compose up -d

down:
	docker-compose down

.PHONY: up down