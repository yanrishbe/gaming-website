run:
	docker run  --rm -e POSTGRES_PASSWORD=docker2147 -e POSTGRES_DB=gaming_website -e POSTGRES_USER=postgres -p 5432:5432 postgres

build:
	POSTGRES_USER=postgres POSTGRES_DB=gaming_website POSTGRES_PASSWORD=docker2147 HOST=localhost PORT=5432 SSLMODE=disable go run ./cmd/main.go
