package connection

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var Connection *sql.DB

func Connect() {

	log.Println(os.Getenv("DATABASE_URL"))
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
	}

	Connection = db

	err = db.Ping()
	if err != nil {
		log.Printf("Unable to connect to database")
		panic(err)
	}

	log.Println("Connected to Postgres DB successfully")
}

func Disconnect() {
	log.Println("Attempting to disconnect from db....")
	err := Connection.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println("Disconnected from db successfully...")
}
