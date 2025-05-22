package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"e-commerce_with_golang/domain"
	"e-commerce_with_golang/domain/mocks"
	"e-commerce_with_golang/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func setUpRouter(productRepo domain.ProductRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	productController := NewProductController(productRepo)

	router.GET("/products", productController.GetProducts)
	router.GET("/products/:id", productController.GetProduct)
	router.POST("/products", productController.CreateProduct)

	return router
}

func TestGetProducts(t *testing.T) {
	mockRepo := mocks.NewProductRepository(t)
	router := setUpRouter(mockRepo)

	mockProducts := []models.Product{
		{
			Model:       gorm.Model{ID: 1},
			Name:        "Product 1",
			Description: "Description 1",
			Price:       10.0,
			Variants:    json.RawMessage(`null`), // Set Variants to a JSON raw message representing null
		},
		{
			Model:       gorm.Model{ID: 2},
			Name:        "Product 2",
			Description: "Description 2",
			Price:       20.0,
			Variants:    json.RawMessage(`null`), // Set Variants to a JSON raw message representing null
		},
	}

	mockRepo.On("GetAll").Return(mockProducts, nil)

	req, _ := http.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Product
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, mockProducts, response)
	mockRepo.AssertExpectations(t)
}

func TestCreateProductWithoutAuth(t *testing.T) {
	mockRepo := mocks.NewProductRepository(t)
	
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	
	router.Use(func(c *gin.Context) {
		c.Set("userID", float64(1))
		c.Next()
	})
	
	productController := NewProductController(mockRepo)
	router.POST("/products", productController.CreateProduct)

	newProduct := models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Category:    "Test Category",
		Quantity:    10,
		Price:       25.99,
		Image:       "test-image.jpg",
		Variants:    json.RawMessage(`[{"color": "red", "image": "red-variant.jpg"}]`),
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Product")).Return(nil).Run(func(args mock.Arguments) {
		product := args.Get(0).(*models.Product)
		product.ID = 1
		t.Logf("Product UserID set to: %d", product.UserID)
	})

	productJSON, _ := json.Marshal(newProduct)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("Response Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response models.Product
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify the product data
	assert.Equal(t, "Test Product", response.Name)
	assert.Equal(t, "Test Description", response.Description)
	assert.Equal(t, "Test Category", response.Category)
	assert.Equal(t, 10, response.Quantity)
	assert.Equal(t, 25.99, response.Price)
	assert.Equal(t, "test-image.jpg", response.Image)
	assert.Equal(t, uint(1), response.UserID) // Should be set by controller
	assert.NotZero(t, response.ID) // Should be set by mock
	
	mockRepo.AssertExpectations(t)
}