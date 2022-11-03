package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/mrityunjaygr8/simplebank/api"
	db "github.com/mrityunjaygr8/simplebank/db/sqlc"
)

const dbDriver = "postgres"
const dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
const address = "0.0.0.0:8080"

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Could not connect to DB")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(address)
	if err != nil {
		log.Fatal("error starting server", err)
	}
}
