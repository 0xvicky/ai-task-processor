package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func Init() (*sql.DB, error) {

	var (
		port              = os.Getenv("PORT")
		POSTGRES_USER     = os.Getenv("POSTGRES_USER")
		POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
		POSTGRES_DB       = os.Getenv("POSTGRES_DB")
	)

	//connection to db
	psqlInfo := fmt.Sprintf("port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		port, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB)

	db, dbErr := sql.Open("postgres", psqlInfo)
	if dbErr != nil {
		return nil, dbErr
	}

	dbPingErr := db.Ping()
	if dbPingErr != nil {
		return nil, dbPingErr
	}
	log.Println("Db connected ✅")

	return db, nil
}
