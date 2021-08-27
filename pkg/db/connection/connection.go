package connection

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gitlab.com/pbobby001/bobpos_api/pkg/logger"
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
	logger.Logger.Info("Attempting to disconnect from db....")
	err := Connection.Close()
	if err != nil {
		_ = logger.Logger.Error(err)
	}
	logger.Logger.Info("Disconnected from db successfully...")
}
