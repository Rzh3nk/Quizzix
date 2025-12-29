<template>
  <div class="quiz-detail-page">
    <div class="background"></div>
    <Header/>
    <div class="content-wrapper">
      <!-- Хлебные крошки -->
      <div class="breadcrumbs">
        <router-link to="/categories" class="breadcrumb-link">
          <span class="breadcrumb-icon">📚</span>
          <span>Категории</span>
        </router-link>
        <span class="breadcrumb-separator">›</span>
        <router-link 
          v-if="quiz?.category" 
          :to="`/category/${quiz.category_id}/quizzes`" 
          class="breadcrumb-link"
        >
          {{ quiz.category.name }}
        </router-link>
        <span v-else class="breadcrumb-current">Загрузка...</span>
        <span class="breadcrumb-separator">›</span>
        <span class="breadcrumb-current">{{ quiz?.title || 'Квиз' }}</span>
      </div>

      
      <div class="quiz-detail-container">
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>Загрузка квиза...</p>
        </div>

        <div v-else-if="error" class="error-state">
          <span class="error-icon">⚠️</span>
          <h3>Квиз не найден</h3>
          <p>{{ error }}</p>
          <button @click="goBack" class="back-btn">
            Вернуться назад
          </button>
        </div>

        <!-- Контент квиза -->
        <div v-else-if="quiz" class="quiz-content">
          <div class="quiz-header">
            <div class="quiz-image-container">
              <div v-if="quiz.img || quiz.ImgURL" class="quiz-image-wrapper">
                <img 
                  :src="getImageUrl(quiz.img || quiz.ImgURL)" 
                  :alt="quiz.img || quiz.ImgURL"
                  class="quiz-image"
                  @error="handleImageError"
                />
              </div>
              <div v-else class="quiz-image-placeholder">
                {{ getQuizIcon(quiz.title) }}
              </div>
            </div>

            <!-- Основная информация -->
            <div class="quiz-info">
              <div class="quiz-meta">
                <span class="quiz-category">
                  <span class="category-icon">🏷️</span>
                  {{ quiz.category?.name || 'Без категории' }}
                </span>
                <span class="quiz-created">
                  <span class="date-icon">📅</span>
                  {{ formatDate(quiz.CreatedAt) }}
                </span>
              </div>

              <h1 class="quiz-title">{{ quiz.title }}</h1>
              
              <div class="quiz-description">
                {{ quiz.description || 'Нет описания' }}
              </div>

              <div class="quiz-stats">
                <div class="stat">
                  <div class="stat-icon">❓</div>
                  <div class="stat-content">
                    <div class="stat-value">{{ quiz.questionCount || 0}}</div>
                    <div class="stat-label">вопросов</div>
                  </div>
                </div>
                
                <div class="stat">
                  <div class="stat-icon">👤</div>
                  <div class="stat-content">
                    <div class="stat-value">{{ formatNumber(quiz.plays || 0) }}</div>
                    <div class="stat-label">прохождений</div>
                  </div>
                </div>
              </div>

              <!-- Автор -->
              <div v-if="quiz.authorName" class="quiz-author">
                <div class="author-avatar">
                  <span class="avatar-icon">👤</span>
                </div>
                <div class="author-info">
                  <div class="author-name">Автор: {{ quiz.authorName }}</div>
                  <div v-if="quiz.authorStats" class="author-stats">
                    {{ quiz.authorStats.quizCount }} квизов · {{ formatNumber(quiz.authorStats.totalPlays) }} прохождений
                  </div>
                </div>
              </div>
            </div>
          </div>

        
          <div class="action-buttons">
            <button @click="startQuiz" class="start-btn">
              <span class="btn-icon">🚀</span>
              <span class="btn-text">Начать квиз</span>
            </button>
            
            <button @click="goBack" class="back-btn">
              <span class="btn-icon">←</span>
              <span class="btn-text">Назад</span>
            </button>

              <button v-if="isAdminUser" 
                @click="deleteQuiz" 
                class="delete-btn">
                <span class="btn-icon">🗑️</span>
                <span class="btn-text">Удалить</span>
              </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth.js'  

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
import Header from '@/components/Header.vue'

const quiz = ref(null)
const loading = ref(true)
const error = ref(null)

const quizId = computed(() => route.params.id)


const fetchQuiz = async () => {
  try {
    loading.value = true
    error.value = null
    
    console.log('Загрузка квиза:', quizId.value)
    
    const quizResponse = await fetch(`/api/quizzes/${quizId.value}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (!quizResponse.ok) {
      if (quizResponse.status === 404) {
        throw new Error('Квиз не найден')
      }
      throw new Error('Ошибка загрузки квиза')
    }
    
    let quizData = await quizResponse.json()
    
    const quizIdNum = quizData.id || quizData.ID
    
    //Количество вопросов
    try {
      const questionsResponse = await fetch(`/api/quizzes/${quizIdNum}/questions`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      })
      if (questionsResponse.ok) {
        const questions = await questionsResponse.json()
        quizData.questionCount = Array.isArray(questions) ? questions.length : 0
      }
    } catch (err) {
      console.warn('Не удалось загрузить вопросы:', err)
    }
    
    //Количество прохождений
    try {
      const resultsResponse = await fetch(`/api/results/${quizIdNum}`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      })
      if (resultsResponse.ok) {
        const resultsData = await resultsResponse.json()
        quizData.plays = resultsData.plays || 0
      }
    } catch (err) {
      console.warn('Не удалось загрузить статистику:', err)
    }
    if (quizData.author_id || quizData.AuthorID) {
      const authorId = quizData.author_id || quizData.AuthorID
      try {
        const authorResponse = await fetch(`/api/users/${authorId}`, {
          headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
        })
        
        if (authorResponse.ok) {
          const author = await authorResponse.json()
          quizData.authorName = author.username || 'Неизвестный автор'
        }
      } catch {
        quizData.authorName = 'Неизвестный автор'
      }
    }
    //Категория
    if (quizData.category_id || quizData.CategoryID) {
      const categoryId = quizData.category_id || quizData.CategoryID
      try {
        const categoryResponse = await fetch(`/api/categories/${categoryId}`, {
          headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
        })
        if (categoryResponse.ok) {
          quizData.category = await categoryResponse.json()
        }
      } catch (err) {
        console.warn('Не удалось загрузить категорию:', err)
      }
    }
    
    quiz.value = quizData
    console.log('✅ Полный квиз:', quiz.value)
    
    
    
  } catch (err) {
    console.error('Ошибка загрузки квиза:', err)
    error.value = err.message || 'Не удалось загрузить квиз'
    
  
  } finally {
    loading.value = false
  }
}


const isAdminUser = computed(() => {
 return authStore.isAdmin.value || localStorage.getItem('role') === 'admin'
})

//Удаление
const deleteQuiz = async () => {
  if (!confirm(`Удалить квиз "${quiz.value.title}"?`)) return
  
  const userId = parseInt(localStorage.getItem('user_id'))
  const quizId = quiz.value.id || quiz.value.ID
  
  const requestBody = {
    user_id: userId,
    quiz_id: quizId
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
    
    alert('✅ Квиз удален!')
    router.push('/categories')  
  } catch (err) {
    alert('❌ ' + err.message)
  }
}

const getImageUrl = (path) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `${API_URL}${path.startsWith('/') ? path : '/' + path}`
}

const handleImageError = (event) => {
  event.target.style.display = 'none'
  event.target.parentElement.classList.add('has-placeholder')
}

const getQuizIcon = (title) => {
  const icons = {
    'история': '📜',
    'наука': '🔬',
    'география': '🌍',
    'литература': '📚',
    'музыка': '🎵',
    'кино': '🎬',
    'спорт': '⚽',
    'технологии': '💻',
    'программирование': '👨‍💻',
    'математика': '🧮'
  }
  
  const lowerTitle = title.toLowerCase()
  for (const [keyword, icon] of Object.entries(icons)) {
    if (lowerTitle.includes(keyword)) return icon
  }
  return '❓'
}


const formatDate = (dateString) => {
  if (!dateString || dateString === 'null' || dateString === 'undefined') {
    return 'Недавно'
  }
  
  const date = new Date(dateString)
  
  if (isNaN(date.getTime())) {
    console.warn('Невалидная дата:', dateString)
    return 'Неизвестно'
  }
  
  return date.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

const formatNumber = (num) => {
  return new Intl.NumberFormat('ru-RU').format(num)
}

const startQuiz = () => {
  router.push(`/quiz/${quizId.value}/play`)
}

const goBack = () => {
  if (quiz.value?.category?.id || quiz.value?.category_id) {
    const categoryId = quiz.value.category?.id || quiz.value.category_id
    router.push(`/category/${categoryId}/quizzes`)
  } else {
    if (window.history.length > 1) {
      router.back()
    } else {
      router.push('/categories')
    }
  }
}

// Наблюдатель за изменением ID квиза
import { watch } from 'vue'

watch(quizId, (newId) => {
  if (newId) {
    fetchQuiz()
  }
})

onMounted(() => {
  if (quizId.value) {
    fetchQuiz()
  }
})
</script>

<style scoped>
@import '@/assets/quiz.css';
</style>