package main

import (
	"log"

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
	r.GET("/api/categories", getCategories)
	r.GET("/api/quizzes/:id/questions", getQuestions)
	r.GET("/api/question/:id/answers", getAnswers)
	r.GET("/api/users/:id", getUser)
	r.POST("/api/register", register)
	r.POST("/api/login", login)
	r.POST("/api/quizzes/:id/submit", checkTest)
	r.GET("/api/users/:id/results", getUserResults)
	r.GET("/api/leaderboard", getLeaderboard)
	r.POST("/api/quizzes", createQuiz)
	log.Println("Backend запущен")

	r.Run(":8080")
}

// тестовые данные
func seedDemoData() {
	var count int64
	db.Model(&models.Quiz{}).Count(&count) //считает количество квизов в бд
	if count > 0 {
		return
	}

	db.Create(&models.Quiz{Title: "Наука"})
	db.Create(&models.Quiz{Title: "Кино"})
	db.Create(&models.Quiz{Title: "История"})

	db.Create(&models.Category{Name: "Наука"})

	db.Create(&models.Question{Text: "Вопрос 1", QuizID: 1})
	db.Create(&models.Question{Text: "Вопрос 2", QuizID: 1})
	db.Create(&models.Question{Text: "Вопрос 3", QuizID: 1})

	db.Create(&models.Answer{Text: "Правильный", QuestionID: 1, IsCorrect: true})
	db.Create(&models.Answer{Text: "Неправильный", QuestionID: 1, IsCorrect: false})

	db.Create(&models.Answer{Text: "Неправильный", QuestionID: 2, IsCorrect: false})
	db.Create(&models.Answer{Text: "Правильный", QuestionID: 2, IsCorrect: true})
}
