package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Rzh3nk/quizzix/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("SUPER_SECRET_KEY") // потом надо скрыть

// структуры для запросов
type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AdminRequest struct {
	CallerID uint `json:"caller_id" binding:"required"`
	TargetID uint `json:"target_id" binding:"required"`
}

func register(c *gin.Context) {
	log.Println("registering user")
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	//сделать более хорошую проверку, когда определим требования по регистрации
	//проверку на уникальность логина, почты, маску на почту, длину всех полей
	if req.Email == "" || req.Password == "" || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing fields"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
		return
	}
	//создание пользователя
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
	}
	//запись в бд
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user exists or db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registered"})
}

func getAllUsers(c *gin.Context) {
	var users []struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}

	if err := db.Model(&models.User{}).
		Select("id, username, email, role").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func login(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	var user models.User
	//проверка по нику(можно сделать еще почту)
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	//создание жвт
	if envSecret := os.Getenv("JWT_SECRET"); envSecret != "" {
		jwtSecret = []byte(envSecret)
	}
	//создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}
	log.Println("login suc")
	//возвращаем данные залогиненного пользователя
	c.JSON(http.StatusOK, gin.H{
		"token":    tokenString,
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	})
}
func setAdmin(c *gin.Context) {
	var req AdminRequest

	// 1. ✅ Оба ID из BODY
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	// 2. ✅ Проверяем CALLER = админ?
	var caller models.User
	if err := db.First(&caller, req.CallerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "caller not found"})
		return
	}
	var target models.User
	if err := db.First(&target, req.TargetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "target not found"})
		return
	}

	if caller.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":       "admin only",
			"caller_id":   req.CallerID,
			"caller_role": caller.Role,
		})
		return
	}
	if target.Role == "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":       "alreafy admin",
			"targer_id":   req.TargetID,
			"target_role": target.Role,
		})
		return
	}
	if req.CallerID == req.TargetID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot promote yourself"})
		return
	}

	result := db.Model(&models.User{}).
		Where("id = ?", req.TargetID).
		Update("role", "admin")

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "target user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "admin role granted",
		"caller_id": req.CallerID,
		"target_id": req.TargetID,
		"success":   true,
	})
}
func getUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"username":    user.Username,
		"email":       user.Email,
		"total_score": user.TotalScore,
		"created":     user.CreatedAt,
		"role":        user.Role,
	})
}
func getUserResults(c *gin.Context) {
	id := c.Param("id")
	var results []models.Result
	if err := db.Where("user_id = ?", id).
		Order("created_at DESC").Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, results)
}
