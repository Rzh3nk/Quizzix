package main

import (
	"log"
	"net/http"

	"github.com/Rzh3nk/quizzix/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("quizzix.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("SQLite не подключен:", err)
	}
	if err := db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Quiz{},
		&models.Question{},
		&models.Answer{},
		&models.Result{},
	); err != nil {
		log.Fatal("Ошибка", err)
	}
	seedDemoData()

	//раутеры
	r := gin.Default()
	r.GET("/api/quizzes", getQuizzes)
	r.GET("/api/quizzes/:id/questions", getQuestions)

	log.Println("Backend запущен")

	r.Run(":8080")
}

func getQuizzes(c *gin.Context) {
	var quizzes []models.Quiz
	if err := db.Find(&quizzes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, quizzes)
}

func getQuestions(c *gin.Context) {
	quizID := c.Param("id")
	var questions []models.Question
	if err := db.Where("quiz_id = ?", quizID).Preload("Answers").Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, questions)
}

// тестовые данные
func seedDemoData() {
	var count int64
	db.Model(&models.Quiz{}).Count(&count) //считает количество квизов в бд
	if count > 0 {
		return
	}

	db.Create(&models.Quiz{Title: "Наука"})
	db.Create(&models.Question{Text: "Вопрос 1", QuizID: 1})
	db.Create(&models.Quiz{Title: "Кино"})
	db.Create(&models.Quiz{Title: "История"})
}
