package main

import (
	"log"

	"github.com/Rzh3nk/quizzix/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
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
	r.GET("/api/quizzes/:id", getQuizByID)
	r.GET("/api/users/:id/quizzes", getQuizzesByAuthor)
	r.GET("/api/categories/:id/quizzes", getCategoryQuizzes)
	r.GET("/api/categories", getCategories)
	r.GET("/api/categories/:id", getCategoryByID)
	r.GET("/api/quizzes/:id/questions", getQuestions)
	r.GET("/api/question/:id/answers", getAnswers)
	r.GET("/api/users/:id", getUser)
	r.POST("/api/register", register)
	r.POST("/api/login", login)
	r.POST("/api/quizzes/:id/submit", checkTest)
	r.GET("/api/users/:id/results", getUserResults)
	r.GET("/api/results/:id", getResultsQuizId)
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
	testhash, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)

	db.Create(&models.User{Username: "test", Email: "test@mail.ru", Password: string(testhash)})
	db.Create(&models.Category{Name: "Наука", Description: "В этой категории представленые науные квизы", ImgURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRpHNKGMwFfIP-QDg_IL1rph5c4Pa5jLos4uA&s"})
	db.Create(&models.Category{Name: "Кино", Description: "В этой категории представленые квизы по кино", ImgURL: "https://1gai.ru/uploads/posts/2021-08/1628577646_1.png"})
	db.Create(&models.Category{Name: "История", Description: "В этой категории представленые исторические квизы", ImgURL: "https://impulse35.ru/wp-content/uploads/2023/11/depositphotos_2478371-stock-photo-old-magnifying-glass-on-word.webp"})

	db.Create(&models.Quiz{Title: "Квиз по науке", Description: "Описание квиза по науке", CategoryID: 1, ImgURL: "https://img.championat.com/news/big/c/x/formula-1-chto-eto-za-sport_17582327151452255140.jpg"})
	db.Create(&models.Quiz{Title: "Квиз по кино", Description: "Описание квиза по кино", CategoryID: 2, ImgURL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRt746EWWaQ2LSMRlfzNNn2f3FTiAU8Vz6jeg&s"})
	db.Create(&models.Quiz{Title: "Квиз по истории", Description: "Описание квиза по истории", CategoryID: 3, ImgURL: "https://auto.vercity.ru/img/magazine/2019/11/11/1573504833.jpg"})

	db.Create(&models.Question{Text: "Вопрос 1", QuizID: 1})
	db.Create(&models.Question{Text: "Вопрос 2", QuizID: 1})
	db.Create(&models.Question{Text: "Вопрос 3", QuizID: 1})

	db.Create(&models.Answer{Text: "Правильный", QuestionID: 1, IsCorrect: true})
	db.Create(&models.Answer{Text: "Неправильный", QuestionID: 1, IsCorrect: false})

	db.Create(&models.Answer{Text: "Неправильный", QuestionID: 2, IsCorrect: false})
	db.Create(&models.Answer{Text: "Правильный", QuestionID: 2, IsCorrect: true})
}
