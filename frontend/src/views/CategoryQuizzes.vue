<template>
  <div class="category-quizzes-page">
    <!-- Фон -->
    <div class="background"></div>
<Header/>
    <!-- Контент -->
    <div class="content-wrapper">
      <!-- Хлебные крошки -->
      <div class="breadcrumbs">
        <router-link to="/categories" class="breadcrumb-link">
          <span class="breadcrumb-icon">📚</span>
          <span>Категории</span>
        </router-link>
        <span class="breadcrumb-separator">›</span>
        <span class="breadcrumb-current">{{ category?.name || 'Загрузка...' }}</span>
      </div>

      <!-- Заголовок категории -->
      <div v-if="category" class="category-header">
        <div v-if="category.imageURL" class="category-header-image">
          <img :src="getImageUrl(category.imageURL)" :alt="category.name" />
        </div>
        <div v-else class="category-header-icon" :style="{ background: getCategoryGradient(category.name) }">
          {{ getCategoryIcon(category.name) }}
        </div>
        <div class="category-header-content">
          <h1>{{ category.name }}</h1>
          <p class="category-description">{{ category.description || 'Описание отсутствует' }}</p>
          <div class="category-stats">
            <div class="stat">
              <span class="stat-icon">📊</span>
              <span class="stat-value">{{ quizzes.length }}</span>
              <span class="stat-label">квизов</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Фильтры -->
      <div class="controls">
        <div class="search-box">
          <span class="search-icon">🔍</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Поиск квизов..."
            @input="handleSearch"
          />
        </div>
        
        <div class="filters">
          <select v-model="selectedDifficulty" @change="applyFilters" class="filter-select">
            <option value="">Все уровни</option>
            <option value="easy">Лёгкий</option>
            <option value="medium">Средний</option>
            <option value="hard">Сложный</option>
          </select>
          
          <select v-model="sortBy" @change="applyFilters" class="filter-select">
            <option value="createdAt">Новые</option>
            <option value="rating">По рейтингу</option>
            <option value="popular">По популярности</option>
          </select>
        </div>
      </div>

      <!-- Лоадер -->
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>Загрузка квизов...</p>
      </div>

      <!-- Ошибка -->
      <div v-else-if="error" class="error-state">
        <span class="error-icon">⚠️</span>
        <p>{{ error }}</p>
        <button @click="fetchCategoryQuizzes" class="retry-btn">
          Попробовать снова
        </button>
      </div>

      <!-- Квизы -->
      <div v-else class="quizzes-container">
        <!-- Счетчик -->
        <div v-if="quizzes.length > 0" class="quizzes-count">
          <span class="count-number">{{ filteredQuizzes.length }}</span>
          <span class="count-text">квизов найдено</span>
        </div>

        <!-- Сетка квизов -->
        <div v-if="filteredQuizzes.length > 0" class="quizzes-grid">
          <div
            v-for="quiz in filteredQuizzes"
            :key="quiz.id"
            class="quiz-card"
            @click="goToQuiz(quiz.ID)"
          >
            <!-- Изображение квиза -->
            <div class="quiz-image-container">
              <img
                v-if="quiz.imageURL"
                :src="getImageUrl(quiz.imageURL)"
                :alt="quiz.title"
                class="quiz-image"
              />
              <div v-else class="quiz-image-placeholder">
                {{ getQuizIcon(quiz.title) }}
              </div>
              <div v-if="quiz.difficulty" class="quiz-difficulty" :class="quiz.difficulty">
                {{ getDifficultyText(quiz.difficulty) }}
              </div>
            </div>

            <!-- Контент квиза -->
            <div class="quiz-content">
              <h3 class="quiz-title">{{ quiz.title }}</h3>
              <p class="quiz-description">{{ quiz.description || 'Без описания' }}</p>
              
              <!-- Статистика квиза -->
              <div class="quiz-stats">
                <div class="quiz-stat">
                  <span class="stat-icon">❓</span>
                  <span class="stat-value">{{ quiz.questionCount || 0 }}</span>
                  <span class="stat-label">вопросов</span>
                </div>
                <div v-if="quiz.timeLimit" class="quiz-stat">
                  <span class="stat-icon">⏱️</span>
                  <span class="stat-value">{{ quiz.timeLimit }}</span>
                  <span class="stat-label">мин</span>
                </div>
                <div v-if="quiz.rating" class="quiz-stat">
                  <span class="stat-icon">⭐</span>
                  <span class="stat-value">{{ quiz.rating.toFixed(1) }}</span>
                </div>
                <div v-if="quiz.plays" class="quiz-stat">
                  <span class="stat-icon">👤</span>
                  <span class="stat-value">{{ quiz.plays }}</span>
                </div>
              </div>

              <!-- Автор -->
              <div v-if="quiz.authorName" class="quiz-author">
                <span class="author-icon">👤</span>
                <span class="author-name">{{ quiz.authorName }}</span>
              </div>
            </div>

            <!-- Кнопка -->
            <button class="start-quiz-btn">
              <span>Начать квиз</span>
              <span class="arrow">→</span>
            </button>
          </div>
        </div>

        <!-- Сообщение если квизов нет -->
        <div v-if="!loading && filteredQuizzes.length === 0" class="empty-state">
          <span class="empty-icon">📭</span>
          <h3>Квизы не найдены</h3>
          <p v-if="searchQuery || selectedDifficulty">
            Попробуйте изменить параметры поиска
          </p>
          <p v-else>
            В этой категории пока нет квизов. Будьте первым!
          </p>
          <button v-if="isAuthenticated" @click="goToCreateQuiz" class="cta-btn">
            Создать квиз
          </button>
          <router-link v-else to="/login" class="auth-link">
            Войдите, чтобы создать квиз
          </router-link>
        </div>
      </div>

      <!-- Кнопка "Назад" -->
      <button @click="goBack" class="back-btn">
        <span class="back-icon">←</span>
        <span>Назад к категориям</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import Header from '@/components/Header.vue'
const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

//const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Данные
const category = ref(null)
const quizzes = ref([])
const loading = ref(false)
const error = ref(null)

// Фильтры
const searchQuery = ref('')
const selectedDifficulty = ref('')
const sortBy = ref('createdAt')

// Получение ID категории из URL
const categoryId = computed(() => route.params.id)

// Загрузка квизов категории
const fetchCategoryQuizzes = async () => {
  try {
    loading.value = true
    error.value = null
    
    console.log(`Загрузка квизов для категории ${categoryId.value}...`)
        const categoryResponse = await fetch(`/api/categories/${categoryId.value}`, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (!categoryResponse.ok) {
      throw new Error(`Категория не найдена: ${categoryResponse.status}`)
    }
    
    category.value = await categoryResponse.json()


    // Формируем URL с параметрами
    let url = `/api/categories/${categoryId.value}/quizzes`
    const params = []
    
    if (selectedDifficulty.value) {
      params.push(`difficulty=${selectedDifficulty.value}`)
    }
    
    if (sortBy.value) {
      params.push(`sort=${sortBy.value}`)
    }
    
    if (searchQuery.value) {
      params.push(`search=${encodeURIComponent(searchQuery.value)}`)
    }
    
    if (params.length > 0) {
      url += `?${params.join('&')}`
    }
    
    // Запрашиваем данные
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`HTTP ${response.status}: ${errorText}`)
    }
    
    const data = await response.json()
    
    quizzes.value = data
    const quizzesWithQuestionCount = await Promise.all(
      (data.quizzes || data).map(async (quiz) => {
        try {
          const quizId = quiz.id || quiz.ID
          if (!quizId) return quiz
          
          // Запрос вопросов
          const questionsResponse = await fetch(`/api/quizzes/${quizId}/questions`, {
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('token')}`
            }
          })
          
          if (questionsResponse.ok) {
            const questions = await questionsResponse.json()
            const questionCount = Array.isArray(questions) ? questions.length : 0
            
            // ✅ Добавляем questionCount
            return {
              ...quiz,
              questionCount: questionCount
            }
          }
        } catch (err) {
          console.warn(`Не удалось загрузить вопросы для квиза ${quiz.id || quiz.ID}:`, err)
        }
        
        return quiz  // Возвращаем без questionCount
      })
    )
    
    // 4. ✅ Сохраняем
    quizzes.value = quizzesWithQuestionCount
    console.log('Данные загружены:', data)
    
  } catch (err) {
    console.error('Ошибка загрузки:', err)
    error.value = err.message || 'Не удалось загрузить квизы'
    
  } finally {
    loading.value = false
  }
}

// Получение URL изображения
const getImageUrl = (imagePath) => {
  if (!imagePath) return ''
  
  if (imagePath.startsWith('http://') || imagePath.startsWith('https://')) {
    return imagePath
  }
  
  if (imagePath.startsWith('/')) {
    return `${API_URL}${imagePath}`
  }
  
  return `${API_URL}/uploads/${imagePath}`
}

// Фильтрация квизов на клиенте (дополнительная к серверной)
const filteredQuizzes = computed(() => {
  let filtered = [...quizzes.value]
  
  // Дополнительный поиск на клиенте (если нужно)
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(quiz =>
      quiz.title.toLowerCase().includes(query) ||
      (quiz.description && quiz.description.toLowerCase().includes(query))
    )
  }
  
  // Дополнительная фильтрация по сложности на клиенте
  if (selectedDifficulty.value) {
    filtered = filtered.filter(quiz => quiz.difficulty === selectedDifficulty.value)
  }
  
  // Сортировка на клиенте
  filtered.sort((a, b) => {
    switch (sortBy.value) {
      case 'rating':
        return (b.rating || 0) - (a.rating || 0)
      case 'popular':
        return (b.plays || 0) - (a.plays || 0)
      default:
        return new Date(b.createdAt) - new Date(a.createdAt)
    }
  })
  
  return filtered
})

// Проверка авторизации
const isAuthenticated = computed(() => {
  return authStore.isAuthenticated?.value || false
})

// Иконки и градиенты
const getCategoryIcon = (categoryName) => {
  const iconMap = {
    'Наука': '🔬',
    'Science': '🔬',
    'Кино': '🎬',
    'Movies': '🎬',
    'История': '📜',
    'History': '📜',
    'География': '🌍',
    'Geography': '🌍',
    'Музыка': '🎵',
    'Music': '🎵',
    'Спорт': '⚽',
    'Sports': '⚽',
    'Литература': '📚',
    'Literature': '📚',
    'Технологии': '💻',
    'Technology': '💻'
  }
  return iconMap[categoryName] || '📁'
}

const getCategoryGradient = (categoryName) => {
  const gradientMap = {
    'Наука': 'linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%)',
    'Science': 'linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%)',
    'Кино': 'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)',
    'Movies': 'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)',
    'История': 'linear-gradient(135deg, #10b981 0%, #34d399 100%)',
    'History': 'linear-gradient(135deg, #10b981 0%, #34d399 100%)',
    'География': 'linear-gradient(135deg, #3b82f6 0%, #60a5fa 100%)',
    'Geography': 'linear-gradient(135deg, #3b82f6 0%, #60a5fa 100%)',
    'Музыка': 'linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)',
    'Music': 'linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)',
    'Спорт': 'linear-gradient(135deg, #ef4444 0%, #f87171 100%)',
    'Sports': 'linear-gradient(135deg, #ef4444 0%, #f87171 100%)'
  }
  return gradientMap[categoryName] || 'linear-gradient(135deg, #6b7280 0%, #9ca3af 100%)'
}

const getQuizIcon = (quizTitle) => {
  return '❓'
}

const getDifficultyText = (difficulty) => {
  const texts = {
    'easy': 'Лёгкий',
    'medium': 'Средний',
    'hard': 'Сложный'
  }
  return texts[difficulty] || difficulty
}

// Методы
const handleSearch = () => {
  // Debounce для поиска
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    applyFilters()
  }, 300)
}

let searchTimeout = null

const applyFilters = () => {
  fetchCategoryQuizzes()
}

const resetFilters = () => {
  searchQuery.value = ''
  selectedDifficulty.value = ''
  sortBy.value = 'createdAt'
  applyFilters()
}

const goToQuiz = (quizId) => {
  router.push(`/quiz/${quizId}`)
}

const goToCreateQuiz = () => {
  if (isAuthenticated.value) {
    router.push(`/create-quiz?category=${categoryId.value}`)
  } else {
    router.push(`/login?redirect=/create-quiz?category=${categoryId.value}`)
  }
}

const goBack = () => {
  router.push('/main')
}

// Наблюдаем за изменением категории в URL
watch(categoryId, () => {
  if (categoryId.value) {
    fetchCategoryQuizzes()
  }
})

// При загрузке страницы
onMounted(() => {
  if (categoryId.value) {
    fetchCategoryQuizzes()
  }
})
</script>

<style scoped>
/* Стили остаются такими же как в предыдущем примере */
/* Можно скопировать стили из предыдущего ответа */
@import '@/assets/catquiz.css';
</style>