package routes_test

import (
	"bytes"
	"e-commerce_with_golang/domain/mocks"
	"e-commerce_with_golang/models"
	"e-commerce_with_golang/routes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() (*gin.Engine, *mocks.ProductRepository, *mocks.ReviewRepository, *mocks.MockUserRepository) {
	// Create mock repositories
	mockProductRepo := new(mocks.ProductRepository)
	mockReviewRepo := new(mocks.ReviewRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	// Setup router with mock dependencies
	router := gin.Default()

	// Simple auth middleware that just checks for "Authorization" header
	authMiddleware := func(c *gin.Context) {
		if c.GetHeader("Authorization") != "Bearer valid-token" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			fmt.Printf("Test aborted dude")
			return
		}
		fmt.Printf("Test Successful dude")
		c.Next()
	}

	routes.SetupRoutes(router, mockProductRepo, mockReviewRepo, authMiddleware)

	return router, mockProductRepo, mockReviewRepo, mockUserRepo
}

func TestGetProductsRoute(t *testing.T) {
	router, mockProductRepo, _, _ := setupTestRouter()

	// Mock the repository response
	mockProducts := []models.Product{
		{
			Model: gorm.Model{ID: 1}, // Using gorm.Model's ID field
			Name:  "Test Product",
			Price: 99.99,
		},
	}
	mockProductRepo.On("GetAll").Return(mockProducts, nil)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/products", nil)

	// Serve request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Product
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, "Test Product", response[0].Name)

	mockProductRepo.AssertExpectations(t)
}

func TestGetProductRoute(t *testing.T) {
	router, mockProductRepo, _, _ := setupTestRouter()

	// Mock the repository response
	mockProduct := &models.Product{
		Model: gorm.Model{ID: 1}, // Using gorm.Model's ID field
		Name:  "Test Product",
		Price: 99.99,
	}
	fmt.Printf("Mock returned: %v", mockProduct)

	mockProductRepo.On("GetByID", uint(1)).Return(mockProduct, nil)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/products/1", nil)

	// Serve request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Product
	err := json.Unmarshal(w.Body.Bytes(), &response)
	fmt.Printf("Response returned: %v", response)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), response.ID) // Checking ID field

	mockProductRepo.AssertExpectations(t)
}

func TestGetReviewsRoute(t *testing.T) {
	router, _, mockReviewRepo, _ := setupTestRouter()

	// Mock the repository response
	mockReviews := []models.Review{
		{
			Model:     gorm.Model{ID: 1}, // Using gorm.Model's ID field
			ProductID: 1,                 // This is the foreign key
			Rating:    5,
			Review:    "Great product",
		},
	}
	mockReviewRepo.On("GetByProductID", uint(1)).Return(mockReviews, nil)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/products/1/reviews", nil)

	// Serve request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Review
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, "Great product", response[0].Review)
	assert.Equal(t, uint(1), response[0].ProductID) // Checking foreign key

	mockReviewRepo.AssertExpectations(t)
}

func TestCreateProductRoute_Unauthorized(t *testing.T) {
	router, _, _, _ := setupTestRouter()

	// Create request without auth header
	w := httptest.NewRecorder()
	productData := models.Product{Name: "New Product", Price: 49.99}
	jsonData, _ := json.Marshal(productData)
	req, _ := http.NewRequest("POST", "/api/products/create", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Serve request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
