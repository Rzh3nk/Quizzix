package models

import "gorm.io/gorm"

//возможно типы ответов
type Answer struct {
	gorm.Model
	QuestionID uint   `gorm:"index" json:"question_id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}
