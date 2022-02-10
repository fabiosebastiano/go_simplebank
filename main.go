package main

import (
	"database/sql"
	"log"

	api "github.com/fabiosebastiano/simplebank/api"
	db "github.com/fabiosebastiano/simplebank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	address  = "localhost:8001"
	dbDriver = "postgres"
	dbSource = "postgresql://root:mysecret@localhost:5432/simple_bank?sslmode=disable"
)

func main() {
	//istanzio DB da passare
	var err error
	connection, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(connection)
	server := api.NewServer(store)
	err = server.Start(address)
	if err != nil {
		log.Fatal("Cannot start the server ", err)
	}
}
