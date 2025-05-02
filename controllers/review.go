package controllers

import (
    "net/http"
    "strconv"
    "e-commerce_with_golang/domain"
    "e-commerce_with_golang/models"
    "github.com/gin-gonic/gin"
)

type ReviewController struct {
    reviewRepo domain.ReviewRepository
}

func NewReviewController(repo domain.ReviewRepository) *ReviewController {
    return &ReviewController{reviewRepo: repo}
}

func (rc *ReviewController) GetReviews(c *gin.Context) {
    productID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    reviews, err := rc.reviewRepo.GetByProductID(uint(productID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, reviews)
}

// Only modify the CreateReview method
func (rc *ReviewController) CreateReview(c *gin.Context) {
    // Get user ID from JWT
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    productID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    var review models.Review
    if err := c.ShouldBindJSON(&review); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Associate review with user and product
    review.ProductID = uint(productID)
    review.UserID = uint(userID.(float64)) // JWT numbers are float64 by default

    if err := rc.reviewRepo.Create(&review); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, review)
}