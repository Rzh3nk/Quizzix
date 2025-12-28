package main

import (
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/Rzh3nk/quizzix/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type quizAnswers struct {
	User_id uint            `json:"user_id"`
	Answers map[uint][]uint `json:"answers"`
}

func checkTest(c *gin.Context) {
	log.Println("test checking")

	var answers quizAnswers
	if err := c.ShouldBindJSON(&answers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	quizID := c.Param("id")
	var questions []models.Question

	// ✅ 1. Загружаем вопросы ПЕРЕД подсчетом
	if err := db.Where("quiz_id = ?", quizID).Preload("Answers").Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	totalQuestions := len(questions) // ✅ ПОСЛЕ Find!
	if totalQuestions == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no questions in quiz"})
		return
	}

	totalScore := 0.0

	// ✅ 2. Подсчет баллов
	for i := 0; i < len(questions); i++ {
		q := questions[i]
		selectedIDs := answers.Answers[q.ID]

		if len(selectedIDs) == 0 {
			continue
		}

		selectedSet := make(map[uint]struct{})
		for j := 0; j < len(selectedIDs); j++ {
			id := selectedIDs[j]
			selectedSet[id] = struct{}{}
		}

		totalCorrect := 0
		correctChosen := 0
		wrongChosen := 0
		correctSet := make(map[uint]struct{})

		for j := 0; j < len(q.Answers); j++ {
			a := q.Answers[j]
			if a.IsCorrect {
				totalCorrect++
				correctSet[a.ID] = struct{}{}
			}
		}

		if totalCorrect == 0 {
			continue
		}

		for id := range selectedSet {
			_, ok := correctSet[id]
			if ok {
				correctChosen++
			} else {
				wrongChosen++
			}
		}

		queScore := float64(correctChosen-wrongChosen) / float64(totalCorrect)
		if queScore < 0 {
			queScore = 0
		}
		totalScore += math.Round(queScore*100) / 100
	}

	// ✅ 3. Парсим quizID
	quizIDint, err := strconv.ParseUint(quizID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quiz id"})
		return
	}
	quizIDuint := uint(quizIDint)

	// ✅ 4. Лучший предыдущий результат (нормальная обработка)
	var bestBefore models.Result
	bestBeforeScore := 0.0
	if err := db.
		Where("user_id = ? AND quiz_id = ?", answers.User_id, quizIDuint).
		Order("score DESC").
		First(&bestBefore).Error; err != nil {

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		// ✅ Нет предыдущих — bestBeforeScore = 0
	} else {
		bestBeforeScore = bestBefore.Score
	}

	// ✅ 5. Обновляем total_score (безопасно)
	if totalScore > bestBeforeScore {
		if err := db.Model(&models.User{}).
			Where("id = ?", answers.User_id).
			Update("total_score", gorm.Expr("COALESCE(total_score, 0) + ?", totalScore-bestBeforeScore)).Error; err != nil {

			log.Printf("Failed to update total_score: %v", err)
			// НЕ возвращаем ошибку — тест прошел успешно
		}
	}

	// ✅ 6. Создаем результат
	result := models.Result{
		Score:    totalScore,
		MaxScore: totalQuestions,
		UserID:   answers.User_id,
		QuizID:   quizIDuint,
	}

	if err := db.Create(&result).Error; err != nil {
		log.Printf("Failed to save result: %v", err)
		// НЕ возвращаем ошибку — тест прошел
	}

	// ✅ 7. Безопасный percent (БЕЗ Inf!)
	percent := 0.0
	if totalQuestions > 0 {
		percent = totalScore / float64(totalQuestions) * 100
	}

	// ✅ 200 OK с правильными данными
	c.JSON(http.StatusOK, gin.H{
		"score":     totalScore,
		"questions": totalQuestions,
		"percent":   percent,
		"max_score": float64(totalQuestions),
		"improved":  totalScore > bestBeforeScore,
	})
}
