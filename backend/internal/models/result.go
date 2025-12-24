package models

import "gorm.io/gorm"

type Result struct {
	gorm.Model
	UserID    uint `gorm:"index" json:"user_id"`
	QuizID    uint `gorm:"index" json:"quiz_id"`
	Score     int  `json:"score"`
	MaxScore  int  `json:"max_score"`
	TimeTaken int  `json:"time_taken"`
}
