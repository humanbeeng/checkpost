run:
	go build -o bin/server && bin/server

create-db:
	docker exec -it db createdb --username=postgres --owner=postgres checkpost

drop-db:
	docker exec -it db dropdb checkpost 

migrate-up:
	migrate -path db/migration -database "postgres://postgres.iwpavmyqhggvhpyjisev:BKsTXPMIBRfpk9Ua@aws-0-us-east-1.pooler.supabase.com:5432/postgres" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://checkpost_owner:UT7ZObH3xufi@ep-wispy-recipe-a43kgz9e.us-east-1.aws.neon.tech/checkpost"  -verbose down

sqlc:
	sqlc generate

.PHONY:
	create-db run drop-db migrate-up migrate-down
