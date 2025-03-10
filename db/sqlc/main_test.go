package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbsource = "postgresql://root:secret@localhost:5432/fintech?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbsource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}
	testQueries = New(conn)
	os.Exit(m.Run())
}
