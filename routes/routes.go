package routes

import (
    "e-commerce_with_golang/controllers"
    "e-commerce_with_golang/domain"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, productRepo domain.ProductRepository, reviewRepo domain.ReviewRepository, authMiddleware gin.HandlerFunc) {
    productController := controllers.NewProductController(productRepo)
    reviewController := controllers.NewReviewController(reviewRepo)

    // Public routes (no authentication needed)
    router.GET("/api/products", productController.GetProducts)
    router.GET("/api/products/:id", productController.GetProduct)
    router.GET("/api/products/:id/reviews", reviewController.GetReviews)

    // Protected routes (require authentication)
    protected := router.Group("/api")
    protected.Use(authMiddleware)
    {
        protected.POST("/products/create", productController.CreateProduct)
        protected.POST("/products/:id/reviews/create", reviewController.CreateReview)
    }
}