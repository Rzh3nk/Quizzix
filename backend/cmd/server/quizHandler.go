package main

import (
	"net/http"

	"github.com/Rzh3nk/quizzix/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func getQuizzes(c *gin.Context) {
	var quizzes []models.Quiz
	if err := db.Find(&quizzes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, quizzes)
}
