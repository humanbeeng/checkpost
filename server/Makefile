run:
	go build -o bin/server && bin/server

create-db:
	docker exec -it db createdb --username=postgres --owner=postgres checkpost

drop-db:
	docker exec -it db dropdb checkpost 

migrate-up:
	migrate -path db/migration -database "postgres://postgres:example@localhost:5432/checkpost?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgres://postgres:example@localhost:5432/checkpost?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY:
	create-db run drop-db migrate-up migrate-down
