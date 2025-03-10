postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

mongodb:
	docker run --name mongodb -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=secret -d mongo:latest

mongosh:
	docker exec -it mongodb bash
	# 		(OR)
	# docker exec -it mongodb mongosh

enterdb:
	mongosh --port 27017 --username root --password secret

dropdb:
	docker exec -it postgres12 dropdb HLS

migratecreate:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/fintech?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/fintech?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown generate test server