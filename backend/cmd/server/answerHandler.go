package main

import (
	"net/http"

	"github.com/Rzh3nk/quizzix/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func getAnswers(c *gin.Context) {
	questionID := c.Param("id")
	var answers []models.Answer
	if err := db.Where("question_id = ?", questionID).Find(&answers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, answers)
}
