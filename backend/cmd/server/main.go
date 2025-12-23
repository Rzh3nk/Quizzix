package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Quiz struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	r := gin.Default()
	r.GET("/api/quizzes", getQuizzes)
	r.Run(":8080")
}

func getQuizzes(c *gin.Context) {
	quizzes := []Quiz{
		{ID: 1, Name: "Наука"},
		{ID: 2, Name: "Кино"},
		{ID: 3, Name: "История"},
	}
	c.JSON(http.StatusOK, quizzes)
}
