package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/thiri-lwin/thiri-bank/api"
	db "github.com/thiri-lwin/thiri-bank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:postgres@localhost:5432/thiri_bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	store := db.NewStore(conn)
	//server := api.NewServer(store)
	//server.Start("0.0.0.0:8080")

	server := api.NewMuxServer(store)
	server.StartMuxServer("0.0.0.0:8080")
}
