package main

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	// ожидание на select sqlite_version()
	mock.ExpectQuery(`select sqlite_version\(\)`).
		WillReturnRows(
			sqlmock.NewRows([]string{"sqlite_version()"}).AddRow("3.40.0"),
		)

	gdb, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	require.NoError(t, err)

	return gdb, mock, sqlDB
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.Default()

	// тут регистрируешь только нужные роуты, либо все сразу
	r.GET("/api/quizzes", getQuizzes)
	r.POST("/api/quizzes", createQuiz)
	r.GET("/api/categories", getCategories)
	r.GET("/api/categories/:id/quizzes", getCategoryQuizzes)
	r.GET("/api/quizzes/:id/questions", getQuestions)
	r.GET("/api/question/:id/answers", getAnswers)
	r.GET("/api/users/:id", getUser)
	r.GET("/api/users/:id/results", getUserResults)
	r.GET("/api/leaderboard", getLeaderboard)
	r.POST("/api/quizzes/:id/submit", checkTest)
	r.POST("/api/register", register)
	r.POST("/api/login", login)

	return r
}
