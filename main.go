package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/mrityunjaygr8/simplebank/api"
	db "github.com/mrityunjaygr8/simplebank/db/sqlc"
	"github.com/mrityunjaygr8/simplebank/utils"
)

func main() {
	config, err := utils.LoadConfig("./.")
	if err != nil {
		log.Fatal("Could not read configuration: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Could not connect to DB")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("error starting server", err)
	}
}
