include .env

air:
	@air -c .air.toml

run:
	@go run cmd/main.go 

migrate-up:
	@migrate -path migration -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?tls=false" -verbose up    

migrate-down:
	@migrate -path migration -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?tls=false" -verbose down

migrate-fix:
	@migrate -path migration -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?tls=false" -verbose force 20241106063649

.PHONY: run migrate-up migrate-down migrate-fix