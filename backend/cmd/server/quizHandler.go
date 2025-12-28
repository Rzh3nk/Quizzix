package main

import (
	"net/http"
	"strconv"

	"github.com/Rzh3nk/quizzix/backend/internal/models"
	"github.com/gin-gonic/gin"
)

type createQuizRequest struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	CategoryID  uint                 `json:"category_id"`
	AuthorID    uint                 `json:"author_id"`
	Questions   []createQuizQuestion `json:"questions"`
}

type createQuizQuestion struct {
	Text    string             `json:"text"`
	Answers []createQuizAnswer `json:"answers"`
}

type createQuizAnswer struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

func getQuizzes(c *gin.Context) {
	var quizzes []models.Quiz
	if err := db.Find(&quizzes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, quizzes)
}

func getQuizByID(c *gin.Context) {
	id := c.Param("id")
	var quiz models.Quiz
	if err := db.Where("id = ?", id).First(&quiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, quiz)
}
func getQuizzesByAuthor(c *gin.Context) {
	authorIDStr := c.Param("id")
	authorID, err := strconv.ParseUint(authorIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid author id"})
		return
	}

	var quizzes []models.Quiz
	if err := db.Where("author_id = ?", uint(authorID)).
		Preload("Category"). // Только категорию preload
		Find(&quizzes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"author_id": uint(authorID),
		"quizzes":   quizzes, // Каждый quiz имеет author_id
		"total":     len(quizzes),
	})
}
func createQuiz(c *gin.Context) {
	var req createQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if req.Title == "" || len(req.Questions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and questions are required"})
		return
	}

	quiz := models.Quiz{
		Title:       req.Title,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		AuthorID:    req.AuthorID, // ✅ Текущий пользователь
	}

	if err := db.Create(&quiz).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// создаём вопросы и ответы
	for i := range req.Questions {
		qReq := req.Questions[i]

		question := models.Question{
			Text:   qReq.Text,
			QuizID: quiz.ID,
		}

		if err := db.Create(&question).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		for j := range qReq.Answers {
			aReq := qReq.Answers[j]

			answer := models.Answer{
				Text:       aReq.Text,
				QuestionID: question.ID,
				IsCorrect:  aReq.IsCorrect,
			}

			if err := db.Create(&answer).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          quiz.ID,
		"title":       quiz.Title,
		"description": quiz.Description,
		"category_id": quiz.CategoryID,
	})
}
