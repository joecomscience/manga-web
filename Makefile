up:
	docker-compose -f scripts/docker-compose.yml up -d
	docker-compose -f scripts/docker-compose.yml logs -f

down:
	docker-compose -f scripts/docker-compose.yml down
