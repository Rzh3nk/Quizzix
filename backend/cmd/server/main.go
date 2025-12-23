package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/api/quizzes", func(w http.ResponseWriter, r *http.Request) {
		quizzes := []string{"Science", "Movies", "History"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"quizzes": quizzes})
	})
	http.ListenAndServe(":8080", nil)
}
