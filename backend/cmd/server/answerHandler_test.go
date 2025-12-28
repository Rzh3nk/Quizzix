package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetAnswers_OK(t *testing.T) {
	gdb, mock, sqlDB := newMockDB(t)
	defer sqlDB.Close()
	db = gdb

	rows := sqlmock.NewRows([]string{"id", "text", "question_id", "is_correct"}).
		AddRow(1, "Ans1", 1, true).
		AddRow(2, "Ans2", 1, false)

	mock.ExpectQuery("SELECT .* FROM `answers` WHERE question_id = ?").
		WithArgs("1").
		WillReturnRows(rows)

	r := setupTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/question/1/answers", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp []map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Len(t, resp, 2)

}
func TestGetAnswers_DBError(t *testing.T) {
	gdb, mock, sqlDB := newMockDB(t)
	defer sqlDB.Close()
	db = gdb

	// эмулируем ошибку БД
	mock.ExpectQuery("SELECT .* FROM `answers` WHERE question_id = ?").
		WithArgs("1").
		WillReturnError(fmt.Errorf("db error"))

	r := setupTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/question/1/answers", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var resp map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, "db error", resp["error"])
}
