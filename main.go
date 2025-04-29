package main

import (
	"e-commerce_with_golang/config"
	"e-commerce_with_golang/database"
	"e-commerce_with_golang/repositories"
	"e-commerce_with_golang/routes"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	database.ConnectDB(cfg)

	// Run migrations
	if err := database.Migrate(database.DB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Initialize repositories
	productRepo := repositories.NewProductRepository(database.DB)
	reviewRepo := repositories.NewReviewRepository(database.DB)

	// Setup routes
	router := routes.SetupRoutes(productRepo, reviewRepo)

	port := "8080"
	log.Printf("Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
