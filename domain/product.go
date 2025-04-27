package domain

import "e-commerce_with_golang/models"

type ProductRepository interface {
    GetAll() ([]models.Product, error)
    GetByID(id uint) (*models.Product, error)
    Create(product *models.Product) error
}