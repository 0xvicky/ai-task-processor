package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Db *sql.DB

func Init() {
	//loading .env vars
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Printf("Fault loading env:%v", envErr)
		return
	}

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

	var dbErr error
	Db, dbErr = sql.Open("postgres", psqlInfo)
	if dbErr != nil {
		panic(dbErr)
	}

	dbPingErr := Db.Ping()
	if dbPingErr != nil {
		log.Fatal(dbPingErr)
	}
	println("Db Connected !!")
}
