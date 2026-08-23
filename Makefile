DB_URL = postgres://postgres:postgres@localhost:5432/productdb?sslmode=disable

run:
	go run cmd/api/main.go

build:
	go build -o bin/api.exe cmd/api/main.go

migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down 1

migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(NAME)

clean:
	rm -rf bin/
	go clean -cache

help:
	@echo "Available commands:"
	@echo "  make run           - run app"
	@echo "  make build         - build binary"
	@echo "  make migrate-up    - apply migrations"
	@echo "  make migrate-down  - rollback one migration"
	@echo "  make migrate-create NAME=name - create migration"
	@echo "  make clean         - clean build"