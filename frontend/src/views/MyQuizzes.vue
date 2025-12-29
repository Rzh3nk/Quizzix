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
            <button @click.stop="deleteQuiz(quiz)" class="delete-btn">
              🗑️ Удалить
            </button>
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

const deleteQuiz = async (quiz) => {
  if (!confirm(`Удалить квиз "${quiz.title || quiz.Title}"?`)) return
  
  const userId = parseInt(localStorage.getItem('user_id'))
  const quizId = quiz.id || quiz.ID
  
  const requestBody = {
    user_id: userId,   // Кто удаляет
    quiz_id: quizId    // Что удаляем
  }
  
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/quizzes/delete', {
      method: 'POST',
      headers: { 
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(requestBody)
    })
    
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Ошибка удаления')
    }
    
    // ✅ Удаляем из списка
    quizzes.value = quizzes.value.filter(q => (q.id || q.ID) !== quizId)
    alert('✅ Квиз удален!')
    
  } catch (err) {
    alert('❌ ' + err.message)
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
/* Основные стили страницы */
.my-quizzes-page {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  margin: 0;
  padding: 0;
  overflow: auto;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, 
    Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
  min-height: 100vh;
}

/* Фон */
.background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: 
    linear-gradient(135deg, #1e1b4b 0%, #312e81 25%, #4f46e5 50%, #7c3aed 75%, #a78bfa 100%),
    repeating-linear-gradient(
      45deg,
      transparent,
      transparent 20px,
      rgba(255, 255, 255, 0.03) 20px,
      rgba(255, 255, 255, 0.03) 40px
    );
  z-index: -1;
}

/* Центральный контейнер */
.content-wrapper {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 30px 20px;
  max-width: 1200px;
  margin: 0 auto;
  text-align: center;
}

/* Шапка страницы */
.header {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 40px;
  padding: 25px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 25px;
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.2);
  color: white;
  border-radius: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  text-decoration: none;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.btn-icon {
  font-size: 1.2rem;
}

.title-section {
  flex-grow: 1;
  text-align: center;
}

.page-title {
  font-size: 2.2rem;
  color: white;
  margin: 0 0 10px 0;
  font-weight: 700;
  text-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}

.subtitle {
  font-size: 1.1rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  font-weight: 400;
}
.delete-btn {
  margin-top: 16px;
  padding: 8px 16px;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  width: 100%;
}

.delete-btn:hover {
  background: #dc2626;
}
.add-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 30px;
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
}

.add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(16, 185, 129, 0.3);
}

/* Состояние загрузки */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 60px 40px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.2);
  margin: 20px 0;
  width: 100%;
}

.spinner {
  width: 60px;
  height: 60px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: #4f46e5;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-state p {
  color: rgba(255, 255, 255, 0.9);
  font-size: 1.1rem;
  margin: 0;
}

/* Состояние ошибки */
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  padding: 50px 40px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(239, 68, 68, 0.3);
  margin: 20px 0;
  width: 100%;
  max-width: 600px;
}

.error-icon {
  font-size: 3.5rem;
  margin-bottom: 10px;
}

.error-state h3 {
  color: white;
  font-size: 1.5rem;
  margin: 0 0 10px 0;
}

.error-state p {
  color: rgba(255, 255, 255, 0.9);
  margin: 0 0 20px 0;
  text-align: center;
  font-size: 1rem;
}

.retry-btn {
  padding: 12px 30px;
  background: linear-gradient(135deg, #ef4444 0%, #f87171 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.retry-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(239, 68, 68, 0.3);
}

/* Состояние "пусто" */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 80px 40px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  backdrop-filter: blur(10px);
  border: 2px dashed rgba(255, 255, 255, 0.3);
  margin: 40px 0;
  width: 100%;
}

.empty-icon {
  font-size: 5rem;
  opacity: 0.7;
}

.empty-state h3 {
  color: white;
  font-size: 1.8rem;
  margin: 0;
}

.empty-state p {
  color: rgba(255, 255, 255, 0.9);
  margin: 0 0 20px 0;
  font-size: 1.1rem;
  max-width: 400px;
}

.create-first-btn {
  padding: 14px 40px;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
}

.create-first-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(79, 70, 229, 0.3);
}

/* Сетка квизов */
.quizzes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 25px;
  margin-top: 20px;
  width: 100%;
}

/* Карточка квиза */
.quiz-card {
  background: white;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  cursor: pointer;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.quiz-card:hover {
  transform: translateY(-10px) scale(1.02);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.25);
}

/* Изображение квиза */
.quiz-image {
  height: 180px;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  background-size: cover;
  background-position: center;
  position: relative;
  display: flex;
  align-items: flex-end;
  padding: 0;
}

.play-count {
  position: absolute;
  top: 15px;
  right: 15px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 5px;
  backdrop-filter: blur(5px);
}

.play-icon {
  font-size: 1rem;
}

/* Контент квиза */
.quiz-content {
  padding: 25px;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.quiz-title {
  font-size: 1.4rem;
  color: #1f2937;
  margin: 0 0 12px 0;
  font-weight: 700;
  line-height: 1.3;
}

.quiz-description {
  color: #6b7280;
  font-size: 0.95rem;
  line-height: 1.5;
  margin-bottom: 20px;
  flex-grow: 1;
}

/* Мета-информация */
.quiz-meta {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-top: 15px;
  border-top: 1px solid #e5e7eb;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #6b7280;
  font-size: 0.9rem;
}

.meta-icon {
  font-size: 1.1rem;
  width: 24px;
  text-align: center;
}

/* Адаптивность */
@media (max-width: 1200px) {
  .quizzes-grid {
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  }
}

@media (max-width: 768px) {
  .content-wrapper {
    padding: 20px 15px;
  }
  
  .header {
    flex-direction: column;
    gap: 20px;
    text-align: center;
    padding: 20px;
  }
  
  .page-title {
    font-size: 1.8rem;
  }
  
  .quizzes-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
  
  .quiz-content {
    padding: 20px;
  }
  
  .quiz-title {
    font-size: 1.3rem;
  }
}

@media (max-width: 480px) {
  .content-wrapper {
    padding: 15px 10px;
  }
  
  .header {
    padding: 15px;
  }
  
  .back-btn, .add-btn {
    width: 100%;
    justify-content: center;
  }
  
  .page-title {
    font-size: 1.6rem;
  }
  
  .subtitle {
    font-size: 1rem;
  }
  
  .empty-state, .error-state, .loading-state {
    padding: 40px 20px;
  }
}

@media (max-height: 700px) {
  .content-wrapper {
    padding-top: 20px;
    padding-bottom: 20px;
  }
  
  .header {
    margin-bottom: 20px;
  }
}
</style>
