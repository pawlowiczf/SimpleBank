package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config file:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot open db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.StartServer(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
