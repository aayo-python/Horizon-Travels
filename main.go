package main

import (
	"database/sql"
	"log"

	api "github.com/janto-pee/fintech-platform.git/api"
	db "github.com/janto-pee/fintech-platform.git/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbsource      = "postgresql://root:secret@localhost:5432/fintech?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbsource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)
}
