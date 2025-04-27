package repositories

import (
    "e-commerce_with_golang/domain"
    "e-commerce_with_golang/models"
    "gorm.io/gorm"
)

type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) GetAll() ([]models.Product, error) {
    var products []models.Product
    err := r.db.Preload("Reviews").Find(&products).Error
    return products, err
}

func (r *productRepository) GetByID(id uint) (*models.Product, error) {
    var product models.Product
    err := r.db.Preload("Reviews").First(&product, id).Error
    return &product, err
}

func (r *productRepository) Create(product *models.Product) error {
    return r.db.Create(product).Error
}