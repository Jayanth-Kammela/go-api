package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Jayanth-Kammela/go-api/database"
	"github.com/Jayanth-Kammela/go-api/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")

	// Initialize the MongoDB client
	database.ConnectDB()

	// Initialize routes
	router := routes.SetupRouter()

	// Start server
	// log.Println("Server is running on port 8080")
	// log.Fatal(http.ListenAndServe(":8080", router))

	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
