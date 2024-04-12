all: webapp

webapp:
	CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static"' -o webapp webapp.go

start-app:
	# Install reflex with 'go install github.com/cespare/reflex@latest'
	# Install godotenv with 'go install github.com/joho/godotenv/cmd/godotenv@latest'
	reflex -s -r '\.go$$' -- godotenv -f .env go run webapp.go

start-view:
	# Install reflex with 'go install github.com/cespare/reflex@latest'
	reflex -r '\.qtpl$$' -- qtc -dir=internal/view

db-migrate:
	migrate -path migrations -database "postgres://127.0.0.1/webapp?sslmode=disable" up

db-schema-dump:
	pg_dump --schema-only -O webapp > internal/database/schema.sql

sqlc-gen:
	sqlc --experimental generate

.PHONY: webapp start-app start-view db-migrate db-schema-dump sqlc-gen
