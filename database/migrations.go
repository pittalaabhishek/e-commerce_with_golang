package database

import (
	"e-commerce_with_golang/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Product{}, &models.Review{}, &models.User{}); err != nil {
		return err
	}
	return nil
}
