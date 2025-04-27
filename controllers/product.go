package controllers

import (
	"net/http"
	"strconv"
	"encoding/json"
	"e-commerce_with_golang/models"
	"e-commerce_with_golang/database"

	"github.com/gin-gonic/gin"
)

type ProductController struct{}

func (pc *ProductController) GetProducts(c *gin.Context) {
	var products []models.Product
	if err := database.DB.Preload("Reviews").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (pc *ProductController) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := database.DB.Preload("Reviews").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert variants to JSONB-compatible format
	variantsJSON, err := json.Marshal(product.Variants)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process variants"})
		return
	}

	// Create a new product with properly formatted variants
	newProduct := models.Product{
		Name:        product.Name,
		Description: product.Description,
		Category:    product.Category,
		Quantity:    product.Quantity,
		Price:       product.Price,
		Image:       product.Image,
		Variants:    variantsJSON,
	}

	if err := database.DB.Create(&newProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newProduct)
}