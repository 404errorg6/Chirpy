package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	dbURL := os.Getenv("DB_URL")
	cfg.secret = os.Getenv("SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error loading database: %v", err)
	}

	cfg.db = database.New(db)
	mux := Mux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", dir, port)
	log.Fatal(server.ListenAndServe())
}
