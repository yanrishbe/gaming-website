run:
	docker run  --rm -e POSTGRES_PASSWORD=docker2147 -e POSTGRES_DB=gaming_website -e POSTGRES_USER=postgres -p 5432:5432 postgres