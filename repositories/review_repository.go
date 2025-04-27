package repositories

import (
    "e-commerce_with_golang/domain"
    "e-commerce_with_golang/models"
    "gorm.io/gorm"
)

type reviewRepository struct {
    db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) domain.ReviewRepository {
    return &reviewRepository{db: db}
}

func (r *reviewRepository) GetByProductID(productID uint) ([]models.Review, error) {
    var reviews []models.Review
    err := r.db.Where("product_id = ?", productID).Find(&reviews).Error
    return reviews, err
}

func (r *reviewRepository) Create(review *models.Review) error {
    if review.Rating < 1 || review.Rating > 5 {
        return models.ErrInvalidRating
    }
    return r.db.Create(review).Error
}