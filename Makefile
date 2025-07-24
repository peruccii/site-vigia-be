# Database connection string - adjust as needed
DATABASE_URL ?= postgres://username:password@localhost/dbname?sslmode=disable

# Migration commands
.PHONY: migrate-up migrate-down migrate-drop migrate-version migrate-force migrate-create sqlc-generate

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-drop:
	migrate -path migrations -database "$(DATABASE_URL)" drop

migrate-version:
	migrate -path migrations -database "$(DATABASE_URL)" version

migrate-force:
	migrate -path migrations -database "$(DATABASE_URL)" force $(VERSION)

migrate-create:
	migrate create -ext sql -dir migrations -seq $(NAME)

sqlc-generate:
	sqlc generate

# Development workflow
.PHONY: dev-setup dev-reset

dev-setup: migrate-up sqlc-generate

dev-reset: migrate-drop migrate-up sqlc-generate
