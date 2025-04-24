build:
	tar -czh . | docker build -t relay - # need to tar because of the symlink
.PHONY: build

up:
	docker compose up -d

logs:
	docker compose logs --tail=100 -f

down ps:
	docker compose $@
