package main

import (
	"net/http"

	"github.com/Rzh3nk/quizzix/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func getQuestions(c *gin.Context) {
	quizID := c.Param("id")
	var questions []models.Question
	if err := db.Where("quiz_id = ?", quizID).Preload("Answers").Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, questions)
}
