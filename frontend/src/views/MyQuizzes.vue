<template>
  <div class="my-quizzes-page">
    <div class="background"></div>
    <Header/>
    <div class="content-wrapper">
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

      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Загрузка ваших квизов...</p>
      </div>

      <div v-else-if="error" class="error-state">
        <span class="error-icon">⚠️</span>
        <h3>Ошибка загрузки</h3>
        <p>{{ error }}</p>
        <button @click="fetchMyQuizzes" class="retry-btn">Повторить</button>
      </div>

      <!--Если пусто-->
      <div v-else-if="quizzes.length === 0" class="empty-state">
        <div class="empty-icon">📝</div>
        <h3>У вас нет квизов</h3>
        <p>Создайте свой первый квиз!</p>
        <router-link to="/create-quiz" class="create-first-btn">
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

const quizzes = ref([])
const loading = ref(true)
const error = ref(null)
const username = ref('Гость')
const userId = ref(0)

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
    
    const quizzesResponse = await fetch(`/api/users/${id}/quizzes`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    
    if (!quizzesResponse.ok) throw new Error('Ошибка загрузки квизов')
    const quizzesData = await quizzesResponse.json()
    let quizzesWithStats = quizzesData.quizzes || []
   
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


onMounted(() => {
  fetchMyQuizzes()
})
</script>

<style scoped>
@import '@/assets/myquiz.css';
</style>
