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
	ImgURL      string               `json:"img"`
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
type DeleteQuizRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	QuizID uint `json:"quiz_id" binding:"required"`
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

	if req.Title == "" || req.ImgURL == "" || len(req.Questions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title, image and questions are required"})
		return
	}

	quiz := models.Quiz{
		Title:       req.Title,
		Description: req.Description,
		ImgURL:      req.ImgURL,
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

func deleteQuiz(c *gin.Context) {
	var req DeleteQuizRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var quiz models.Quiz
	if err := db.First(&quiz, req.QuizID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
		return
	}

	var user models.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if quiz.AuthorID != req.UserID && user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":           "not authorized",
			"quiz_author_id":  quiz.AuthorID,
			"request_user_id": req.UserID,
			"user_role":       user.Role,
		})
		return
	}

	if err := db.Delete(&models.Quiz{}, req.QuizID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "quiz deleted",
		"quiz_id": req.QuizID,
		"user_id": req.UserID,
		"success": true,
	})
}
