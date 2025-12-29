<template>
  <div class="profile-page">
    <!-- Фон -->
    <div class="background"></div>
    <Header/>
    
    <!-- Контент -->
    <div class="content-wrapper">
      <!-- Назад -->
      <button @click="goBack" class="back-btn">← Назад</button>

      <!-- Загрузка -->
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>Загрузка профиля...</p>
      </div>

      <!-- Ошибка -->
      <div v-else-if="error" class="error">
        <span class="error-icon">⚠️</span>
        <h3>Ошибка загрузки</h3>
        <p>{{ error }}</p>
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

<style scoped>
/* Основные стили страницы */
.profile-page {
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
  max-width: 800px;
  margin: 0 auto;
  text-align: center;
}

/* Кнопка назад */
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
  margin-bottom: 30px;
  align-self: flex-start;
  text-decoration: none;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

/* Состояние загрузки */
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 20px;
  padding: 60px 40px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.2);
  width: 100%;
  color: rgba(255, 255, 255, 0.9);
  font-size: 1.2rem;
}

/* Состояние ошибки */
.error {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 50px 40px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(239, 68, 68, 0.3);
  width: 100%;
  color: white;
  text-align: center;
}

.error button {
  padding: 10px 25px;
  background: linear-gradient(135deg, #ef4444 0%, #f87171 100%);
  color: white;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 10px;
}

.error button:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(239, 68, 68, 0.3);
}

/* Основной контейнер профиля */
.profile {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 30px;
}

/* Информация о профиле */
.profile-info {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 40px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.2);
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.profile-info h1 {
  font-size: 2.5rem;
  color: white;
  margin: 0;
  font-weight: 700;
  text-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}

.profile-info p {
  font-size: 1.2rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  display: flex;
  align-items: center;
  gap: 10px;
}

.score {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 15px 30px;
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
  border-radius: 15px;
  color: white;
  font-size: 1.3rem;
  font-weight: 700;
  margin-top: 10px;
}

.score-icon {
  font-size: 1.5rem;
}

/* История прохождений */
.history-section {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 30px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.history-section h2 {
  font-size: 1.8rem;
  color: white;
  margin: 0 0 25px 0;
  text-align: center;
  font-weight: 600;
}

.no-results {
  color: rgba(255, 255, 255, 0.7);
  font-size: 1.1rem;
  text-align: center;
  padding: 40px 20px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 15px;
  border: 2px dashed rgba(255, 255, 255, 0.2);
}

/* Список результатов */
.results-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.result-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 15px;
  transition: all 0.3s ease;
  border: 2px solid rgba(255, 255, 255, 0.1);
}

.result-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
  border-color: rgba(79, 70, 229, 0.3);
}

.quiz-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  text-align: left;
}

.quiz-info strong {
  font-size: 1.1rem;
  color: #1f2937;
  font-weight: 600;
}

.quiz-info span {
  font-size: 0.9rem;
  color: #6b7280;
}

.result-score {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  padding: 8px 15px;
  border-radius: 10px;
  font-weight: 600;
  font-size: 1rem;
  min-width: 140px;
  text-align: center;
}

/* Действия профиля */
.profile-actions {
  display: flex;
  gap: 20px;
  justify-content: center;
  flex-wrap: wrap;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 15px 35px;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  border-radius: 12px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.3s ease;
  min-width: 180px;
}

.action-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 15px 35px rgba(79, 70, 229, 0.3);
}

/* Адаптивность */
@media (max-width: 768px) {
  .content-wrapper {
    padding: 20px 15px;
  }
  
  .profile-info {
    padding: 30px 20px;
  }
  
  .profile-info h1 {
    font-size: 2rem;
  }
  
  .history-section {
    padding: 25px 20px;
  }
  
  .history-section h2 {
    font-size: 1.6rem;
  }
  
  .result-item {
    flex-direction: column;
    gap: 15px;
    text-align: center;
    padding: 20px;
  }
  
  .quiz-info {
    text-align: center;
  }
  
  .result-score {
    width: 100%;
  }
  
  .profile-actions {
    flex-direction: column;
  }
  
  .action-btn {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .content-wrapper {
    padding: 15px 10px;
  }
  
  .profile-info {
    padding: 25px 15px;
  }
  
  .profile-info h1 {
    font-size: 1.8rem;
  }
  
  .profile-info p {
    font-size: 1rem;
  }
  
  .score {
    padding: 12px 25px;
    font-size: 1.1rem;
  }
  
  .history-section {
    padding: 20px 15px;
  }
  
  .history-section h2 {
    font-size: 1.4rem;
  }
  
  .back-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>