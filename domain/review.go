package domain

import "e-commerce_with_golang/models"

type ReviewRepository interface {
    GetByProductID(productID uint) ([]models.Review, error)
    Create(review *models.Review) error
}