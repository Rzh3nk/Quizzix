package main

import (
	"net/http"

	"github.com/Rzh3nk/quizzix/backend/internal/models"
	"github.com/gin-gonic/gin"
)

func getCategories(c *gin.Context) {
	var categories []models.Category
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, categories)
}
func getCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, category)
}

func getCategoryQuizzes(c *gin.Context) {
	categoryID := c.Param("id")
	var quizzes []models.Quiz
	if err := db.Where("category_id = ?", categoryID).Find(&quizzes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, quizzes)
}
