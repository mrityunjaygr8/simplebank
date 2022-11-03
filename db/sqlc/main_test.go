package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/mrityunjaygr8/simplebank/utils"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("./../..")
	if err != nil {
		log.Fatal("Count not load the configuration: ", err)
	}
	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Could not connect to DB:", err)
	}

	testQueries = New(testDb)
	os.Exit(m.Run())
}
