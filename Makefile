.PHONY: generate

migrate-up:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/boilerplate-1?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/boilerplate-1?sslmode=disable" -verbose down