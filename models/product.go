package models

import (
    "encoding/json"
    "gorm.io/gorm"
)

type Variant struct {
    Color string `json:"color"`
    Image string `json:"image"`
}

type Product struct {
    gorm.Model
    Name        string          `json:"name" gorm:"not null"`
    Description string          `json:"description"`
    Category    string          `json:"category"`
    Quantity    int             `json:"quantity" gorm:"default:0"`
    Price       float64         `json:"price" gorm:"not null"`
    Image       string          `json:"image"`
    Variants    json.RawMessage `json:"variants" gorm:"type:jsonb"`
    Reviews     []Review        `json:"reviews" gorm:"foreignKey:ProductID"`
}