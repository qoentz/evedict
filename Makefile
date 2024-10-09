run:
	go run cmd/app/main.go

migrate:
	go run cmd/migrate/main.go

create-migration:
	migrate create -ext sql -dir internal/db/migrations ${name}