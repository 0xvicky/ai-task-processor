package main

import (
	"ai-task-processor/internal/handler"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	println("AI-TASK-PROCESSOR")
	//loading .env vars
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal(envErr)
	}

	var (
		port              = os.Getenv("PORT")
		POSTGRES_USER     = os.Getenv("POSTGRES_USER")
		POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
		POSTGRES_DB       = os.Getenv("POSTGRES_DB")
	)
	//Register Routes
	http.HandleFunc("/", handler.RootHandler)
	http.HandleFunc("/health", handler.Health)

	//connection to db
	psqlInfo := fmt.Sprintf("port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		port, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB)

	db, dbErr := sql.Open("postgres", psqlInfo)
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	dbPingErr := db.Ping()
	if dbPingErr != nil {
		log.Fatal(dbPingErr)
	}
	println("Db Connected !!")

	//Server
	err := http.ListenAndServe(":6969", nil)

	if err != nil {
		log.Fatal(err)
	}

	// println("Server running at 6969")
}
