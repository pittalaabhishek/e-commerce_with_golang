package models

import "gorm.io/gorm"

type Variant struct {
	Color string `json:"color"`
	Image string `json:"image"`
}

type Product struct {
	gorm.Model
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Image       string    `json:"image"`
	Variants    []Variant `json:"variants" gorm:"type:jsonb"`
	Reviews     []Review  `json:"reviews" gorm:"foreignKey:ProductID"`
}