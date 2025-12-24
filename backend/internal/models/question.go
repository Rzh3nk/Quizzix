package models

import "gorm.io/gorm"

//добавить типы вопросов и типы ответов (несколько ответов, вопрос с картинкой и подобное)
type Question struct {
	gorm.Model
	QuizID    uint   `gorm:"index" json:"quiz_id"`
	Text      string `json:"text"`
	TimeLimit int    `json:"time_limit"`

	Answers []Answer `gorm:"foreignKey:QuestionID" json:"answers"`
}
