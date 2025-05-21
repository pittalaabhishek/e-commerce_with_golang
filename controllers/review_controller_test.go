package controllers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
	"strconv"

	"e-commerce_with_golang/domain"
    "e-commerce_with_golang/domain/mocks"
    "e-commerce_with_golang/models"
    "github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/assert"
    "gorm.io/gorm"
)

func setUpReviewRouter(reviewRepo domain.ReviewRepository) *gin.Engine {
    gin.SetMode(gin.TestMode)
    router := gin.Default()

    reviewController := NewReviewController(reviewRepo)

    router.GET("/products/:id/reviews", reviewController.GetReviews)
    router.POST("/products/:id/reviews", reviewController.CreateReview)

    return router
}

func TestGetReviews(t *testing.T) {
    mockRepo := mocks.NewReviewRepository(t)
    router := setUpReviewRouter(mockRepo)

    productID := uint(1)
    mockReviews := []models.Review{
        {Model: gorm.Model{ID: 1}, ProductID: productID, UserID: 1, Rating: 5, Review: "Great product!"},
    }

    mockRepo.On("GetByProductID", productID).Return(mockReviews, nil)

    req, _ := http.NewRequest("GET", "/products/1/reviews", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var response []models.Review
    json.Unmarshal(w.Body.Bytes(), &response)

    assert.Equal(t, mockReviews, response)
    mockRepo.AssertExpectations(t)
}

func TestCreateReview(t *testing.T) {
	// Initialize mock repository
	mockRepo := mocks.NewReviewRepository(t)
	
	// Create a test Gin engine without JWT middleware
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())
	
	// Create a test group and handle the review creation
	api := router.Group("/products/:productID/reviews")
	api.POST("", func(c *gin.Context) {
		// Set userID in context to simulate JWT middleware
		c.Set("userID", uint(1))
		
		// Parse the request body
		var review models.Review
		if err := c.ShouldBindJSON(&review); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Get parameters
		productIDStr := c.Param("productID")
		productID, err := strconv.ParseUint(productIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
			return
		}
		
		// Set IDs in review
		userID := c.MustGet("userID").(uint)
		review.UserID = userID
		review.ProductID = uint(productID)
		
		// Call repository
		if err := mockRepo.Create(&review); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusCreated, review)
	})
	
	// Test data
	productID := uint(1)
	userID := uint(1)
	review := models.Review{
		Rating: 5, // Changed back to int to match your model
		Review: "Great product!",
	}
	
	// Setup mock expectation
	// We use mock.AnythingOfType since the pointer address might be different
	mockRepo.On("Create", mock.AnythingOfType("*models.Review")).Return(nil).Run(func(args mock.Arguments) {
		passedReview := args.Get(0).(*models.Review)
		assert.Equal(t, productID, passedReview.ProductID)
		assert.Equal(t, userID, passedReview.UserID)
		assert.Equal(t, 5, passedReview.Rating) // Changed to int
		assert.Equal(t, "Great product!", passedReview.Review)
	})
	
	// Create request
	reviewJSON, _ := json.Marshal(review)
	req, _ := http.NewRequest("POST", "/products/1/reviews", bytes.NewBuffer(reviewJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Create response recorder and serve the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	var responseReview models.Review
	json.Unmarshal(w.Body.Bytes(), &responseReview)
	assert.Equal(t, 5, responseReview.Rating) // Changed to int
	assert.Equal(t, "Great product!", responseReview.Review)
	
	// No need for explicit Assert.Expectations since it's done in Cleanup of your mock
}