package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"peruccii/site-vigia-be/internal/api"
	"peruccii/site-vigia-be/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func connectWithRetry(databaseURL string, maxRetries int) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			log.Printf("Attempt %d: Failed to open database connection: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("Attempt %d: Failed to ping database: %v", i+1, err)
			db.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Successfully connected to database")
		return db, nil
	}

	return nil, err
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	log.Printf("Using DATABASE_URL: %s", databaseURL)

	db, err := connectWithRetry(databaseURL, 10)
	if err != nil {
		log.Fatal("Failed to connect to database after retries:", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	router := api.SetupRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}