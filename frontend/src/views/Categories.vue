<template>
  <div class="categories-page"><Header/>
    <div class="background"></div>
    
    <div class="content-wrapper">
      
      <div class="logo-row">
        <div class="logo-circle">Q</div>
        <h1>Категории Quizzix</h1>
      </div>

      <p class="subtitle">Выберите интересующую категорию и проверьте свои знания</p>

      <!--Поиск-->
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

      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <p>Загрузка категорий...</p>
      </div>

      <div v-else-if="error" class="error-state">
        <span class="error-icon">⚠️</span>
        <p>{{ error }}</p>
        <button @click="fetchCategories" class="retry-btn">
          Попробовать снова
        </button>
      </div>

      <div v-else class="categories-container">
        <!--Счетчик категрий-->
        <div class="categories-count">
          <span class="count-text">Найдено категорий: </span>
          <span class="count-number">{{ filteredCategories.length }}</span>
        </div>

        <!--Сам список категорий-->
        <div class="categories-grid">
          <div
            v-for="category in filteredCategories"
            :key="category.ID"
            class="category-card"
            @click="goToCategory(category.ID)"
          >
            <!--Верхняя часть карточки категории-->
            <div class="category-header" :style="{ background: getCategoryGradient(category.name) }">
              <img 
                v-if="category.img"
                :src="category.img"
                alt="quiz image"
                class="quiz-image"
              />
            </div>

            <div class="category-content">
              <h3 class="category-title">{{ category.name }}</h3>
              <p class="category-description">{{ category.description || `Квизы по теме "${category.name}"` }}</p>

              <div class="category-stats">
                <div class="stat-item">
                  <span class="stat-icon">📊</span>
                  <span class="stat-label">Квизов: </span>
                  <span class="stat-value">{{ category.QuizCount || 0 }}</span>
                </div>
              </div>
            </div>

            <button class="explore-btn">
              <span>Начать</span>
              <span class="arrow">→</span>
            </button>
          </div>
        </div>

        <div v-if="filteredCategories.length === 0 && categories.length > 0" class="empty-state">
          <span class="empty-icon">🔍</span>
          <h3>Категории не найдены</h3>
          <p>Попробуйте изменить поисковый запрос</p>
          <button @click="clearSearch" class="clear-btn">
            Очистить поиск
          </button>
        </div>
      </div>

      <div v-if="isAuthenticated" class="cta-section">
        <div class="cta-card">
          <h3>Хотите создать свой квиз?</h3>
          <p>Создавайте увлекательные тесты и делитесь ими с другими</p>
          <button @click="goToCreateQuiz" class="cta-btn">
            Создать квиз
          </button>
        </div>
      </div>

      <div class="navigation">
        <router-link to="/" class="nav-link">
          <span class="nav-icon">🏠</span>
          <span>Главная</span>
        </router-link>
        <router-link v-if="isAuthenticated" to="/profile" class="nav-link">
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

const searchQuery = ref('')
const loading = ref(false)
const error = ref(null)
const categories = ref([])

//Загрузка категорий с API
const fetchCategories = async () => {
  try {
    loading.value = true
    error.value = null
    
    console.log('Загрузка категорий с API...')
    
    const response = await fetch(`/api/categories`, {
      headers: {
        'Content-Type': 'application/json',
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
    
    categories.value = data.map(category => ({
      ...category,
      QuizCount: category.QuizCount || 0,
    }))
    
    console.log('Категории загружены:', categories.value)
    
  } catch (err) {
    console.error('Ошибка загрузки категорий:', err)
    error.value = err.message || 'Не удалось загрузить категории'
  } finally {
    loading.value = false
  }
}

//Фильтрация категорий по поиску
const filteredCategories = computed(() => {
  if (!searchQuery.value) return categories.value
  
  const query = searchQuery.value.toLowerCase()
  return categories.value.filter(category => 
    category.name.toLowerCase().includes(query) 
  )
})

//Проверка авторизации
const isAuthenticated = computed(() => {
  return authStore.isAuthenticated?.value || false
})



//Градиенты для категорий
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


// Методы
const handleSearch = () => {
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
onMounted(() => {
  fetchCategories()
})
</script>

<style scoped>
@import '@/assets/categories.css';
</style>