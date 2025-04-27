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

func (rc *ReviewController) CreateReview(c *gin.Context) {
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

    review.ProductID = uint(productID)
    if err := rc.reviewRepo.Create(&review); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, review)
}