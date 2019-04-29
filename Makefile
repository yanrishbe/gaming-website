CONN:=user=postgres dbname=gaming_website password=docker2147 host=localhost port=5432 sslmode=disable

run:
	docker run  --rm -e POSTGRES_PASSWORD=docker2147 -e POSTGRES_DB=gaming_website -e POSTGRES_USER=postgres -p 5432:5432 postgres

build:
	export CONN="user=postgres dbname=gaming_website password=docker2147 host=localhost port=5432 sslmode=disable"
	CONN="user=postgres dbname=gaming_website password=docker2147 host=localhost port=5432 sslmode=disable" go run ./cmd/main.go
