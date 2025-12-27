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
	totalScore := 0.0

	quizID := c.Param("id")
	var questions []models.Question
	totalQuestions := len(questions)
	if err := c.ShouldBindJSON(&answers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := db.Where("quiz_id = ?", quizID).Preload("Answers").Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

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
	quizIDint, _ := strconv.ParseUint(quizID, 10, 64)
	quizIDuint := uint(quizIDint)

	var bestBefore models.Result
	bestBeforeScore := 0.0

	if err := db.
		Where("user_id = ? AND quiz_id = ?", answers.User_id, quizIDuint).
		Order("score DESC").
		First(&bestBefore).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			bestBeforeScore = 0.0
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
	} else {
		bestBeforeScore = bestBefore.Score
	}

	if totalScore > bestBeforeScore {
		if err := db.Model(&models.User{}).
			Where("id = ?", answers.User_id).
			Update("total_points", gorm.Expr("total_points + ?", totalScore-bestBeforeScore)).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
	}

	result := models.Result{
		Score:    totalScore,
		MaxScore: totalQuestions,
		UserID:   answers.User_id,
		QuizID:   quizIDuint,
	}

	if err := db.Create(&result).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "result exists or db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"score":     totalScore,
		"questions": totalQuestions,
		"percent":   totalScore / float64(totalQuestions) * 100,
	})
}
