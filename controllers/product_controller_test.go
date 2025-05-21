package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"e-commerce_with_golang/domain"
	"e-commerce_with_golang/domain/mocks"
	"e-commerce_with_golang/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
