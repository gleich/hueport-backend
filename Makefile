reset:
	docker compose down
	docker system prune -af
	docker volume rm hueport-scraper_postgres_data

start:
	docker compose up -d
	sleep 2
	go run cmd/scraper.go