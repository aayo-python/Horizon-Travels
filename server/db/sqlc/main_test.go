package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/janto-pee/fintech-platform.git/util"
	_ "github.com/lib/pq"
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgresql://root:secret@localhost:5432/fintech?sslmode=disable"
// )

var testQueries *Queries
var TestDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	TestDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	testQueries = New(TestDB)
	os.Exit(m.Run())
}
