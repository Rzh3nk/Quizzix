package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex" json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}
