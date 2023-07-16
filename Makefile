up:
	docker-compose build
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs

restart: down up