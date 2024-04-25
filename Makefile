proto-gen:
	./scripts/gen-proto.sh

create-migrate:
	migrate create -ext sql -dir ./migrations -seq users-tables

migrate-up:
	migrate -path ./migrations -database 'postgres://postgres:123@localhost:5432/establishmentdb?sslmode=disable' up

migrate-down:
	migrate -path ./migrations -database 'postgres://postgres:123@localhost:5432/establishmentdb?sslmode=disable' down

migrate-force:
	migrate -path ./migrations -database 'postgres://postgres:123@localhost:5432/establishmentdb?sslmode=disable' force 1
