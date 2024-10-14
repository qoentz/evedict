run:
	go run cmd/app/main.go

migrate:
	go run cmd/migrate/main.go

create-migration:
	migrate create -ext sql -dir internal/db/migrations ${name}

air:
	air -c air.toml

tailwind:
	npx tailwindcss -i internal/view/css/input.css -o static/styles.css --watch

templ:
	templ generate -watch -proxy=http://localhost:8080
