package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"
)

const dbConnStr = "postgres://diceroller:diceroller@localhost/diceroller?sslmode=disable"

func main() {
	// TODO use the flag package to parse dbConnStr and port

	// Seed according to time
	rand.Seed(time.Now().UnixNano())

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal("Error while connecting to postgres: ", err)
	}
	// Sanity check that the DB is connectable
	if err := db.Ping(); err != nil {
		log.Fatal("Error while testing connection to db: ", err)
	}

	s := NewServer(db)
	err = s.Run(":8080")
	if err != nil {
		log.Fatal("Error while running api server: ", err)
	}
}
