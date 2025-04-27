package main

import (
	"log"
	"os"
	"e-commerce_with_golang/database"
	"e-commerce_with_golang/config"
	"e-commerce_with_golang/migrations"
	"e-commerce_with_golang/routes"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.LoadConfig()

	// Initialize database
	database.ConnectDB(cfg)
	db := database.DB

	// Run migrations
	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Setup routes
	router := routes.SetupRoutes()

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // default port
	}
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}