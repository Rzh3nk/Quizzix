package models

import "gorm.io/gorm"

type Quiz struct {
	gorm.Model
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ImgURL      string     `json:"img"`
	CategoryID  uint       `gorm:"index" json:"category_id"`
	Category    Category   `json:"-"`
	AuthorID    uint       `gorm:"index" json:"author_id"`
	Questions   []Question `gorm:"foreignKey:QuizID" json:"-"`
	Results     []Result   `gorm:"foreignKey:QuizID" json:"-"`
}
