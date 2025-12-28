package main

import (
	"net/http"
	"strconv"

	"github.com/Rzh3nk/quizzix/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func getLeaderboard(c *gin.Context) {
	var users []models.User
	if err := db.Order("total_score DESC").Limit(100).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func getResultsQuizId(c *gin.Context) {
	quizID := c.Param("id")
	quizIDUint := 0

	// Парсим ID
	if id64, err := strconv.ParseUint(quizID, 10, 32); err == nil {
		quizIDUint = int(id64)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz id"})
		return
	}

	var count int64
	if err := db.Model(&models.Result{}).
		Where("quiz_id = ?", quizIDUint).
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"quiz_id": quizIDUint,
		"plays":   count, // ✅ Количество прохождений
	})
}
