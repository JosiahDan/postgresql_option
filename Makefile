createdb:
	docker exec -it postgre_local createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgre_local dropdb simple_bank

postgres:
	docker run --name postgre_local -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres

migrateup:
	migrate -path db/migration -database "postgresql://root:123456@192.168.18.130:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:123456@192.168.18.130:5432/simple_bank?sslmode=disable" -verbose down

.PHONY: createdb dropdb postgres migrateup