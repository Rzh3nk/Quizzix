<template>
  <div class="my-quizzes-page">
    <!-- Фон -->
    <div class="background"></div>
    <Header/>
    <!-- Контент -->
    <div class="content-wrapper">
      <!-- Шапка -->
      <div class="header">
        <button @click="goBack" class="back-btn">
          <span class="btn-icon">←</span>
          Назад
        </button>
        <div class="title-section">
          <h1 class="page-title">📚 Мои квизы</h1>
          <p class="subtitle">{{ username }} • {{ quizzes.length }} квизов</p>
        </div>
        <router-link to="/create-quiz" class="add-btn">
          <span class="btn-icon">+</span>
          Создать
        </router-link>
      </div>

      <!-- Загрузка -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Загрузка ваших квизов...</p>
      </div>

      <!-- Ошибка -->
      <div v-else-if="error" class="error-state">
        <span class="error-icon">⚠️</span>
        <h3>Ошибка загрузки</h3>
        <p>{{ error }}</p>
        <button @click="fetchMyQuizzes" class="retry-btn">Повторить</button>
      </div>

      <!-- Пустой список -->
      <div v-else-if="quizzes.length === 0" class="empty-state">
        <div class="empty-icon">📝</div>
        <h3>У вас нет квизов</h3>
        <p>Создайте свой первый квиз!</p>
        <router-link to="/quiz-add" class="create-first-btn">
          Создать первый квиз
        </router-link>
      </div>

      <!-- Список квизов -->
      <div v-else class="quizzes-grid">
        <div 
          v-for="quiz in quizzes" 
          :key="quiz.id || quiz.ID"
          class="quiz-card"
          @click="goToQuiz(quiz)"
        >
          <div class="quiz-image" :style="{ backgroundImage: `url(${quiz.img || quiz.ImgURL})` }">
            <div class="play-count">
              <span class="play-icon">👤</span>
              {{ formatNumber(quiz.plays || 0) }}
            </div>
          </div>
          
          <div class="quiz-content">
            <h3 class="quiz-title">{{ quiz.title || quiz.Title }}</h3>
            <p class="quiz-description">{{ quiz.description || quiz.Description }}</p>
            
            <div class="quiz-meta">
              <div class="meta-item">
                <span class="meta-icon">❓</span>
                {{ getQuestionCount(quiz) }} вопросов
              </div>
              <div class="meta-item">
                <span class="meta-icon">🏷️</span>
                {{ quiz.category?.name || 'Без категории' }}
              </div>
              <div class="meta-item">
                <span class="meta-icon">📅</span>
                {{ formatDate(quiz.created_at || quiz.CreatedAt) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Header from '@/components/Header.vue'
const router = useRouter()

// Данные
const quizzes = ref([])
const loading = ref(true)
const error = ref(null)
const username = ref('Гость')
const userId = ref(0)

// Получить ID пользователя
const getUserId = () => {
  const userIdStr = localStorage.getItem('user_id')
  return userIdStr ? parseInt(userIdStr) : 0
}

// Загрузка квизов
const fetchMyQuizzes = async () => {
  try {
    loading.value = true
    error.value = null
    
    const id = getUserId()
    if (!id) {
      throw new Error('Пользователь не авторизован')
    }
    
    userId.value = id
    console.log('🔍 Загрузка квизов пользователя:', id)
    
    const token = localStorage.getItem('token')
    
    // 1. ✅ Основные квизы
    const quizzesResponse = await fetch(`/api/users/${id}/quizzes`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    
    if (!quizzesResponse.ok) throw new Error('Ошибка загрузки квизов')
    const quizzesData = await quizzesResponse.json()
    let quizzesWithStats = quizzesData.quizzes || []
    
    // 2. ✅ questionCount + plays для КАЖДОГО квиза
    const quizPromises = quizzesWithStats.map(async (quiz) => {
      const quizId = quiz.id || quiz.ID
      
      // Вопросы
      try {
        const questionsRes = await fetch(`/api/quizzes/${quizId}/questions`, {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        if (questionsRes.ok) {
          const questions = await questionsRes.json()
          quiz.questionCount = Array.isArray(questions) ? questions.length : 0
        }
      } catch {}

      // Прохождения
      try {
        const resultsRes = await fetch(`/api/results/${quizId}`, {
          headers: { 'Authorization': `Bearer ${token}` }
        })
        if (resultsRes.ok) {
          const resultsData = await resultsRes.json()
          quiz.plays = resultsData.plays || 0
        }
      } catch {}

      // Категория
      if (quiz.category_id) {
        try {
          const categoryRes = await fetch(`/api/categories/${quiz.category_id}`, {
            headers: { 'Authorization': `Bearer ${token}` }
          })
          if (categoryRes.ok) {
            quiz.category = await categoryRes.json()
          }
        } catch {}
      }
      
      return quiz
    })
    
    quizzesWithStats = await Promise.all(quizPromises)
    quizzes.value = quizzesWithStats
    
    // Username
    const userResponse = await fetch(`/api/users/${id}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (userResponse.ok) {
      const userData = await userResponse.json()
      username.value = userData.username
    }
    
  } catch (err) {
    console.error('Ошибка загрузки квизов:', err)
    error.value = err.message
  } finally {
    loading.value = false
  }
}

// Навигация
const goBack = () => {
  if (window.history.length > 1) {
      router.back()
    } else {
      router.push('/main')
    }
}

const goToQuiz = (quiz) => {
  router.push(`/quiz/${quiz.id || quiz.ID}`)
}

// Вспомогательные функции
const getQuestionCount = (quiz) => {
  return quiz.questionCount || 0
}

const formatNumber = (num) => {
  return new Intl.NumberFormat('ru-RU').format(num)
}

const formatDate = (dateStr) => {
  if (!dateStr) return 'Недавно'
  try {
    const date = new Date(dateStr)
    return date.toLocaleDateString('ru-RU', {
      day: 'numeric',
      month: 'short',
      year: 'numeric'
    })
  } catch {
    return 'Неизвестно'
  }
}

// Загрузка при монтировании
onMounted(() => {
  fetchMyQuizzes()
})
</script>

<style scoped>
.my-quizzes-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

.background {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: 
    radial-gradient(circle at 20% 50%, rgba(255,255,255,0.1) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(255,255,255,0.05) 0%, transparent 50%);
  z-index: -1;
}

.content-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  min-height: 100vh;
}

/* Header */
.header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 40px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(20px);
  padding: 24px 32px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.back-btn, .add-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 28px;
  border-radius: 16px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.back-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  cursor: pointer;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.add-btn {
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
  margin-left: auto;
}

.add-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 15px 35px rgba(16, 185, 129, 0.4);
}

.title-section {
  flex: 1;
  text-align: center;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 800;
  background: linear-gradient(135deg, #fff 0%, #f0f9ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0 0 8px 0;
}

.subtitle {
  color: rgba(255, 255, 255, 0.9);
  font-size: 1.2rem;
  margin: 0;
}

/* Состояния */
.loading-state, .error-state, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 60px 40px;
  text-align: center;
  margin: 40px 0;
}

.spinner {
  width: 64px;
  height: 64px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-top-color: #10b981;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 24px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-icon {
  font-size: 6rem;
  margin-bottom: 24px;
  opacity: 0.7;
}

.retry-btn, .create-first-btn {
  padding: 14px 32px;
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
  border: none;
  border-radius: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
  margin-top: 20px;
}

.retry-btn:hover, .create-first-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 15px 35px rgba(16, 185, 129, 0.4);
}

/* Сетка квизов */
.quizzes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 24px;
  margin-top: 40px;
}

.quiz-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
}

.quiz-card:hover {
  transform: translateY(-12px);
  box-shadow: 0 35px 70px rgba(0, 0, 0, 0.2);
}

.quiz-image {
  height: 200px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  background-size: cover;
  background-position: center;
  position: relative;
  display: flex;
  align-items: flex-end;
}

.play-count {
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(10px);
  color: white;
  padding: 8px 16px;
  border-radius: 20px 20px 0 0;
  margin: 0 16px;
  font-size: 0.9rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
}

.quiz-content {
  padding: 28px;
}

.quiz-title {
  font-size: 1.4rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 12px 0;
  line-height: 1.3;
}

.quiz-description {
  color: #64748b;
  font-size: 0.95rem;
  line-height: 1.5;
  margin: 0 0 20px 0;
}

.quiz-meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
  color: #475569;
}

.meta-icon {
  font-size: 1.1rem;
}

/* Адаптивность */
@media (max-width: 768px) {
  .header {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }
  
  .add-btn {
    margin-left: 0 !important;
  }
  
  .quizzes-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .quiz-content {
    padding: 24px;
  }
}
</style>
