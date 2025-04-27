package routes

import (
    "e-commerce_with_golang/controllers"
    "e-commerce_with_golang/domain"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(productRepo domain.ProductRepository, reviewRepo domain.ReviewRepository) *gin.Engine {
    r := gin.Default()

    productController := controllers.NewProductController(productRepo)
    reviewController := controllers.NewReviewController(reviewRepo)

    // Product routes
    r.GET("/api/products", productController.GetProducts)
    r.GET("/api/products/:id", productController.GetProduct)
    r.POST("/api/products/create", productController.CreateProduct)

    // Review routes
    r.GET("/api/products/:id/reviews", reviewController.GetReviews)
    r.POST("/api/products/:id/reviews/create", reviewController.CreateReview)

    return r
}