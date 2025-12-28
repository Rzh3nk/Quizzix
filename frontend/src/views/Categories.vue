<template>
  <div class="categories-page">

    <!-- Фон -->
    <div class="background"></div>
    <Header/>
    <!-- Контент -->
    <div class="content-wrapper">
      
      <!-- Логотип и название -->
      <div class="logo-row">
        <div class="logo-circle">Q</div>
        <h1>Категории Quizzix</h1>
      </div>

      <!-- Подзаголовок -->
      <p class="subtitle">Выберите интересующую категорию и проверьте свои знания</p>

      <!-- Поиск -->
      <div class="search-container">
        <div class="input-wrapper">
          <span class="input-icon">🔍</span>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Поиск категории..."
            @input="handleSearch"
          />
          <button v-if="searchQuery" @click="clearSearch" class="clear-search-btn">
            ✕
          </button>
        </div>
      </div>

      <!-- Лоадер -->
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>Загрузка категорий...</p>
      </div>

      <!-- Ошибка -->
      <div v-else-if="error" class="error-state">
        <span class="error-icon">⚠️</span>
        <p>{{ error }}</p>
        <button @click="fetchCategories" class="retry-btn">
          Попробовать снова
        </button>
      </div>

      <!-- Категории -->
      <div v-else class="categories-container">
        <!-- Счетчик -->
        <div class="categories-count">
          <span class="count-number">{{ filteredCategories.length }}</span>
          <span class="count-text">категорий найдено</span>
        </div>

        <!-- Сетка категорий -->
        <div class="categories-grid">
          <div
            v-for="category in filteredCategories"
            :key="category.ID"
            class="category-card"
            @click="goToCategory(category.ID)"
          >
            <!-- Верхняя часть с иконкой -->
            <div class="category-header" :style="{ background: getCategoryGradient(category.name) }">
              <img 
            v-if="category.img"
            :src="category.img"
            alt="quiz image"
            class="quiz-image"
          />
              <div v-if="category.isPopular" class="category-badge">🔥</div>
            </div>

            <!-- Контент категории -->
            <div class="category-content">
              <h3 class="category-title">{{ category.name }}</h3>
              <p class="category-description">{{ category.description || `Квизы по теме "${category.name}"` }}</p>

              <!-- Статистика -->
              <div class="category-stats">
                <div class="stat-item">
                  <span class="stat-icon">📊</span>
                  <span class="stat-value">{{ category.QuizCount || 0 }}</span>
                  <span class="stat-label">квизов</span>
                </div>
                <div class="stat-item">
                  <span class="stat-icon">⭐</span>
                  <span class="stat-value">{{ category.Difficulty || 'Средняя' }}</span>
                  <span class="stat-label">сложность</span>
                </div>
              </div>

              <!-- Дата создания -->
              <div v-if="category.CreatedAt" class="category-date">
                <span class="date-icon">📅</span>
                <span class="date-text">{{ formatDate(category.CreatedAt) }}</span>
              </div>
            </div>

            <!-- Кнопка -->
            <button class="explore-btn">
              <span>Начать</span>
              <span class="arrow">→</span>
            </button>
          </div>
        </div>

        <!-- Сообщение если ничего не найдено -->
        <div v-if="filteredCategories.length === 0 && categories.length > 0" class="empty-state">
          <span class="empty-icon">🔍</span>
          <h3>Категории не найдены</h3>
          <p>Попробуйте изменить поисковый запрос</p>
          <button @click="clearSearch" class="clear-btn">
            Очистить поиск
          </button>
        </div>

        <!-- Сообщение если категорий нет вообще -->
        <div v-if="categories.length === 0 && !loading" class="empty-state">
          <span class="empty-icon">📭</span>
          <h3>Категории пока отсутствуют</h3>
          <p>Будьте первым, кто создаст категорию!</p>
          <button v-if="isAuthenticated" @click="goToCreateCategory" class="cta-btn">
            Создать категорию
          </button>
          <router-link v-else to="/login" class="auth-link">
            Войдите, чтобы создать категорию
          </router-link>
        </div>
      </div>

      <!-- Призыв к действию -->
      <div v-if="isAuthenticated" class="cta-section">
        <div class="cta-card">
          <h3>Хотите создать свой квиз?</h3>
          <p>Создавайте увлекательные тесты и делитесь ими с другими</p>
          <button @click="goToCreateQuiz" class="cta-btn">
            Создать квиз
          </button>
        </div>
      </div>

      <!-- Навигация -->
      <div class="navigation">
        <router-link to="/" class="nav-link">
          <span class="nav-icon">🏠</span>
          <span>Главная</span>
        </router-link>
        <router-link v-if="isAuthenticated" to="/dashboard" class="nav-link">
          <span class="nav-icon">📊</span>
          <span>Профиль</span>
        </router-link>
        <router-link v-else to="/login" class="nav-link">
          <span class="nav-icon">🔐</span>
          <span>Войти</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'

import Header from '@/components/Header.vue'

const router = useRouter()
const authStore = useAuthStore()
const quizCounts = ref({})
//const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

const searchQuery = ref('')
const loading = ref(false)
const error = ref(null)
const categories = ref([])

// Загрузка категорий с API
const fetchCategories = async () => {
  try {
    loading.value = true
    error.value = null
    
    console.log('Загрузка категорий с API...')
    
    const response = await fetch(`/api/categories`, {
      headers: {
        'Content-Type': 'application/json',
        //'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }
    
    const data = await response.json()
    await Promise.all(
      data.map(async (category) => {
        try {
          const quizzesRes = await fetch(`/api/categories/${category.ID}/quizzes`)
          const quizzes = await quizzesRes.json()
          category.QuizCount = quizzes.length
        } catch (err) {
          category.QuizCount = 0
        }
      })
    )
    // Добавляем дополнительные поля для отображения
    categories.value = data.map(category => ({
      ...category,
      QuizCount: category.QuizCount || 0,
      Difficulty: category.Difficulty || 'Средняя',
      isPopular: category.QuizCount > 10 // Пример логики для популярных
    }))
    
    console.log('Категории загружены:', categories.value)
    
  } catch (err) {
    console.error('Ошибка загрузки категорий:', err)
    error.value = err.message || 'Не удалось загрузить категории'
    
    // Fallback данные для демонстрации
    if (process.env.NODE_ENV === 'development') {
      categories.value = getMockCategories()
    }
  } finally {
    loading.value = false
  }
}
const loadQuizCount = async (categoryId) => {
  try {
    const response = await fetch(`/api/categories/${categoryId}/quizzes`)
    const quizzes = await response.json()
    quizCounts.value[categoryId] = quizzes.length  // ✅ Считаем!
  } catch (err) {
    quizCounts.value[categoryId] = 0
  }
}
// Фильтрация категорий по поиску
const filteredCategories = computed(() => {
  if (!searchQuery.value) return categories.value
  
  const query = searchQuery.value.toLowerCase()
  return categories.value.filter(category => 
    category.name.toLowerCase().includes(query) 
  )
})

// Проверка авторизации
const isAuthenticated = computed(() => {
  return authStore.isAuthenticated?.value || false
})



// Градиенты для категорий
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
    'Sports': 'linear-gradient(135deg, #ef4444 0%, #f87171 100%)',
    'Литература': 'linear-gradient(135deg, #ec4899 0%, #f472b6 100%)',
    'Literature': 'linear-gradient(135deg, #ec4899 0%, #f472b6 100%)',
    'Технологии': 'linear-gradient(135deg, #06b6d4 0%, #22d3ee 100%)',
    'Technology': 'linear-gradient(135deg, #06b6d4 0%, #22d3ee 100%)'
  }
  return gradientMap[categoryName] || 'linear-gradient(135deg, #6b7280 0%, #9ca3af 100%)'
}

// Форматирование даты
const formatDate = (dateString) => {
  if (!dateString) return ''
  
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('ru-RU', {
      day: 'numeric',
      month: 'long',
      year: 'numeric'
    })
  } catch {
    return dateString
  }
}


// Методы
const handleSearch = () => {
  // Можно добавить debounce для оптимизации
  console.log('Поиск:', searchQuery.value)
}

const clearSearch = () => {
  searchQuery.value = ''
}

const goToCategory = (categoryId) => {
  router.push(`/category/${categoryId}/quizzes`)
}

const goToCreateQuiz = () => {
  if (isAuthenticated.value) {
    router.push('/create-quiz')
  } else {
    router.push('/login?redirect=/create-quiz')
  }
}

const goToCreateCategory = () => {
  router.push('/create-category')
}

// При загрузке страницы
onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
/* Стили только для этого компонента */
@import '@/assets/categories.css';
</style>