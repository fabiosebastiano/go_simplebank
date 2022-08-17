package main

import (
	"database/sql"
	"log"

	api "github.com/fabiosebastiano/simplebank/api"
	db "github.com/fabiosebastiano/simplebank/db/sqlc"
	"github.com/fabiosebastiano/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot read the config file ", err)
	}

	//istanzio DB da passare
	connection, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(connection)
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("Cannot create the server ", err)
	}

	err = server.Start(config.ServerHost + ":" + config.ServerPort)
	if err != nil {
		log.Fatal("Cannot start the server ", err)
	}
}
