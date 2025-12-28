✅ Profile.vue — Профиль пользователя (без стилей)
text
<template>
  <div class="profile-page">
         <Header/>
    <!-- Назад -->
    <button @click="goBack" class="back-btn">← Назад</button>

    <!-- Загрузка -->
    <div v-if="loading" class="loading">Загрузка профиля...</div>

    <!-- Ошибка -->
    <div v-else-if="error" class="error">
      {{ error }}
      <button @click="fetchProfile">Повторить</button>
    </div>

    <!-- Профиль -->
    <div v-else-if="user" class="profile">
      <!-- Основная информация -->
      <div class="profile-info">
        <h1>{{ user.username }}</h1>
        <p>📧 {{ user.email }}</p>
        <div class="score">
          <span class="score-icon">⭐</span>
          {{ formatNumber(user.total_score || 0) }} баллов
        </div>
      </div>

      <!-- История прохождений -->
      <div class="history-section">
        <h2>Последние прохождения</h2>
        <div v-if="results.length === 0" class="no-results">
          Пока нет прохождений
        </div>
        <div v-else class="results-list">
          <div 
            v-for="result in results" 
            :key="result.id"
            class="result-item"
          >
            <div class="quiz-info">
              <strong>{{ result.quiz_title || 'Квиз' }}</strong>
              <span>{{ formatDate(result.created_at) }}</span>
            </div>
            <div class="result-score">
              {{ result.score }} / {{ result.max_score }} 
              ({{ Math.round((result.score / result.max_score) * 100) }}%)
            </div>
          </div>
        </div>
      </div>

      <!-- Навигация -->
      <div class="profile-actions">
        <router-link to="/my-quizzes" class="action-btn">Мои квизы</router-link>
        <router-link to="/leaderboard" class="action-btn">Лидерборд</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>

import Header from '@/components/Header.vue'
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Данные
const user = ref(null)
const results = ref([])
const loading = ref(true)
const error = ref(null)

// ID пользователя
const getUserId = () => {
  const userIdStr = localStorage.getItem('user_id')
  return userIdStr ? parseInt(userIdStr) : 0
}

// Загрузка профиля
const fetchProfile = async () => {
  try {
    loading.value = true
    error.value = null
    
    const id = getUserId()
    if (!id) throw new Error('Пользователь не авторизован')
    
    const token = localStorage.getItem('token')
    
    // 1. Информация о пользователе
    const userRes = await fetch(`/api/users/${id}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (!userRes.ok) throw new Error('Ошибка загрузки пользователя')
    user.value = await userRes.json()
    
    // 2. Последние результаты
    const resultsRes = await fetch(`/api/users/${id}/results`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (resultsRes.ok) {
      const data = await resultsRes.json()
      results.value = data.slice(0, 10) // Топ-10 последних
    }
    
  } catch (err) {
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

const formatNumber = (num) => {
  return new Intl.NumberFormat('ru-RU').format(num)
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Загрузка
onMounted(() => {
  fetchProfile()
})
</script>