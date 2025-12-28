<template>
  <div class="category-quizzes-page">
    <!-- Фон -->
    <div class="background"></div>

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
            @click="goToQuiz(quiz.id)"
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
    
    category.value = data.category
    quizzes.value = data.quizzes || []
    
    console.log('Данные загружены:', data)
    
  } catch (err) {
    console.error('Ошибка загрузки:', err)
    error.value = err.message || 'Не удалось загрузить квизы'
    
    // Mock данные для разработки
    if (process.env.NODE_ENV === 'development') {
      category.value = {
        id: categoryId.value,
        name: 'История',
        description: 'Исторические события, личности и факты',
        imageURL: null
      }
      
      quizzes.value = [
        {
          id: 1,
          title: 'История Древнего Рима',
          description: 'Проверьте свои знания о Римской империи',
          imageURL: null,
          difficulty: 'medium',
          timeLimit: 20,
          questionCount: 15,
          rating: 4.7,
          plays: 1240,
          authorName: 'Иван Иванов',
          createdAt: new Date().toISOString()
        },
        {
          id: 2,
          title: 'Вторая мировая война',
          description: 'Важные события и личности войны',
          imageURL: null,
          difficulty: 'hard',
          timeLimit: 30,
          questionCount: 25,
          rating: 4.8,
          plays: 890,
          authorName: 'Мария Петрова',
          createdAt: new Date().toISOString()
        }
      ]
    }
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
  router.push('/categories')
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

/* Основные контейнеры */
.category-quizzes-page {
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

/* Хлебные крошки */
.breadcrumbs {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 30px;
  width: 100%;
  color: white;
}

.breadcrumb-link {
  display: flex;
  align-items: center;
  gap: 5px;
  color: white;
  text-decoration: none;
  opacity: 0.8;
  transition: opacity 0.3s ease;
}

.breadcrumb-link:hover {
  opacity: 1;
}

/* Заголовок категории */
.category-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 30px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 15px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.2);
  width: 100%;
}

.category-header-image {
  width: 80px;
  height: 80px;
  border-radius: 15px;
  overflow: hidden;
  flex-shrink: 0;
}

.category-header-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.category-header-icon {
  width: 80px;
  height: 80px;
  border-radius: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2.5rem;
  flex-shrink: 0;
  color: white;
}

.category-header-content {
  flex-grow: 1;
  text-align: left;
}

.category-header h1 {
  font-size: 1.8rem;
  margin-bottom: 10px;
  text-align: left;
  color: white;
}

.category-description {
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 15px;
  font-size: 1rem;
  text-align: left;
}

.category-stats {
  display: flex;
  gap: 20px;
}

.stat {
  display: flex;
  align-items: center;
  gap: 5px;
  color: rgba(255, 255, 255, 0.8);
}

.stat-icon {
  font-size: 1rem;
}

/* Управляющие элементы */
.controls {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-bottom: 30px;
  width: 100%;
}

.search-box {
  position: relative;
  width: 100%;
}

.search-icon {
  position: absolute;
  left: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: #6b7280;
}

.search-box input {
  width: 100%;
  padding: 12px 15px 12px 45px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-size: 1rem;
  background: white;
  color: #1f2937;
}

.filters {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.filter-select {
  padding: 10px 15px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  background: white;
  color: #1f2937;
  font-size: 0.95rem;
  cursor: pointer;
  min-width: 150px;
}

/* Лоадер и ошибки */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  padding: 50px 0;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 15px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(239, 68, 68, 0.3);
  max-width: 500px;
  margin: 20px auto;
}

.error-icon {
  font-size: 3rem;
}

/* Сетка квизов */
.quizzes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 25px;
  margin-bottom: 40px;
  width: 100%;
}

.quiz-card {
  background: white;
  border-radius: 15px;
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

.quiz-image-container {
  position: relative;
  height: 150px;
  overflow: hidden;
}

.quiz-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.quiz-card:hover .quiz-image {
  transform: scale(1.05);
}

.quiz-image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 3rem;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
}

.quiz-difficulty {
  position: absolute;
  top: 10px;
  right: 10px;
  padding: 5px 10px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  background: rgba(0, 0, 0, 0.7);
  color: white;
}

.quiz-content {
  padding: 20px;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

.quiz-title {
  font-size: 1.3rem;
  margin-bottom: 10px;
  color: #1f2937;
}

.quiz-description {
  color: #6b7280;
  font-size: 0.95rem;
  margin-bottom: 15px;
  line-height: 1.5;
  min-height: 60px;
}

.quiz-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  margin-bottom: 15px;
}

.quiz-stat {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #6b7280;
  font-size: 0.9rem;
}

.quiz-stat .stat-icon {
  font-size: 1rem;
}

.quiz-stat .stat-value {
  font-weight: 600;
  color: #4f46e5;
}

.quiz-author {
  display: flex;
  align-items: center;
  gap: 5px;
  color: #9ca3af;
  font-size: 0.85rem;
  margin-top: auto;
}

.start-quiz-btn {
  width: calc(100% - 40px);
  margin: 0 20px 20px;
  padding: 12px;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.start-quiz-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(79, 70, 229, 0.3);
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
  margin-top: 30px;
  backdrop-filter: blur(10px);
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

/* Адаптивность */
@media (max-width: 768px) {
  .category-header {
    flex-direction: column;
    text-align: center;
  }
  
  .category-header-content {
    text-align: center;
  }
  
  .filters {
    flex-direction: column;
  }
  
  .filter-select {
    width: 100%;
  }
  
  .quizzes-grid {
    grid-template-columns: 1fr;
  }
}
</style>