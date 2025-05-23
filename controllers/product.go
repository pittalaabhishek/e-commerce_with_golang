package controllers

import (
    "net/http"
    "strconv"
    "e-commerce_with_golang/domain"
    "e-commerce_with_golang/models"
    "github.com/gin-gonic/gin"
)

type ProductController struct {
    productRepo domain.ProductRepository
}

func NewProductController(repo domain.ProductRepository) *ProductController {
    return &ProductController{productRepo: repo}
}

func (pc *ProductController) GetProducts(c *gin.Context) {
    products, err := pc.productRepo.GetAll()
    if err != nil {
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

    product, err := pc.productRepo.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }
    c.JSON(http.StatusOK, product)
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
    // Get user ID from JWT
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Associate product with user
    product.UserID = uint(userID.(float64)) // JWT numbers are float64 by default

    if err := pc.productRepo.Create(&product); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, product)
}