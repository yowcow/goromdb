all: build up logs

build:
	docker-compose build

up:
	docker-compose up -d

logs:
	docker-compose logs -f

stop:
	docker-compose stop

down:
	docker-compose down

.PHONY: all build up logs stop down
