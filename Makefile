up:
	docker-compose up -d

down:
	docker-compose up -d

createdb:
	docker exec -it stuhub-be-db-1 createdb --username=postgres --owner=postgres stuhub

dropdb:
	docker exec -it stuhub-be-db-1 dropdb --username=postgres stuhub

migrateup:
	migrate -path migrations -database "postgresql://postgres:password@localhost:5432/stuhub?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations -database "postgresql://postgres:password@localhost:5432/stuhub?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

.PHONY: down createdb dropdb migrateup migratedown