<template>
  <div class="leaderboard-page">
    <!-- Фон -->
    <div class="background-gradient"></div>
    <Header/>
    <!-- Контент -->
    <div class="content-wrapper">
      <!-- Заголовок -->
      <div class="header">
        <button @click="goBack" class="back-btn">
          <span class="btn-icon">←</span>
          <span class="btn-text">Назад</span>
        </button>
        <div class="title-section">
          <h1 class="page-title">🏆 Лидерборд</h1>
          <p class="subtitle">Топ-100 игроков по очкам</p>
        </div>
      </div>

      <!-- Загрузка -->
      <div v-if="loading" class="loading-container">
        <div class="spinner"></div>
        <p>Загрузка лидерборда...</p>
      </div>

      <!-- Ошибка -->
      <div v-else-if="error" class="error-container">
        <span class="error-icon">⚠️</span>
        <h3>Ошибка загрузки</h3>
        <p>{{ error }}</p>
        <button @click="fetchLeaderboard" class="retry-btn">
          Попробовать снова
        </button>
      </div>

      <!-- Таблица лидеров -->
      <div v-else class="leaderboard-container">
        <div class="leaderboard-table">
          <div class="table-header">
            <div class="header-cell rank">#</div>
            <div class="header-cell name">Игрок</div>
            <div class="header-cell score">Очки</div>
          </div>
          
          <div 
            v-for="(user, index) in leaderboard" 
            :key="user.id || user.ID"
            class="table-row"
            :class="{ 'top-3': index < 3, 'current-user': isCurrentUser(user) }"
          >
            <div class="cell rank">
              <span v-if="index < 3" class="medal">
                {{ ['🥇', '🥈', '🥉'][index] }}
              </span>
              <span v-else>{{ index + 1 }}</span>
            </div>
            <div class="cell name">
              <span class="username">{{ user.username }}</span>
              <span v-if="isCurrentUser(user)" class="current-badge">👑 Вы</span>
            </div>
            <div class="cell score">
              {{ formatNumber(user.total_score || user.total_points || 0) }}
            </div>
          </div>
        </div>

        <!-- Пустой список -->
        <div v-if="leaderboard.length === 0" class="empty-state">
          <span class="empty-icon">👥</span>
          <h3>Пока нет игроков</h3>
          <p>Будьте первым в лидерборде!</p>
        </div>

        <!-- Информация -->
        <div class="info">
          <p>Обновлено: {{ formatDate(updatedAt) }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Header from '@/components/Header.vue'
const router = useRouter()
const route = useRoute()

// Данные
const leaderboard = ref([])
const loading = ref(true)
const error = ref(null)
const updatedAt = ref(new Date())

// Загрузка лидерборда
const fetchLeaderboard = async () => {
  try {
    loading.value = true
    error.value = null
    
    const token = localStorage.getItem('token')
    const response = await fetch('/api/leaderboard', {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (!response.ok) {
      throw new Error(`Ошибка ${response.status}`)
    }
    
    const data = await response.json()
    leaderboard.value = Array.isArray(data) ? data : data.users || []
    updatedAt.value = new Date()
    
  } catch (err) {
    console.error('Leaderboard error:', err)
    error.value = err.message || 'Не удалось загрузить лидерборд'
  } finally {
    loading.value = false
  }
}

// Текущий пользователь
const currentUserId = computed(() => {
  try {
    const userStr = localStorage.getItem('user') || localStorage.getItem('user_id')
    if (!userStr) return null
    
    const user = typeof userStr === 'string' ? parseInt(userStr) : JSON.parse(userStr)
    return user?.id || user?.user_id || null
  } catch {
    return null
  }
})

const isCurrentUser = (user) => {
  const userId = user.id || user.ID || user.user_id
  return currentUserId.value === userId
}

// Форматирование
const formatNumber = (num) => {
  return new Intl.NumberFormat('ru-RU').format(num)
}

const formatDate = (date) => {
  return new Date().toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Навигация
const goBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/categories')
  }
}

// Загрузка при монтировании
onMounted(() => {
  fetchLeaderboard()
})
</script>

<style scoped>
.leaderboard-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #1e1b4b 0%, #312e81 50%, #4f46e5 100%);
  position: relative;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

.background-gradient {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 20% 80%, rgba(120, 119, 198, 0.3) 0%, transparent 50%),
    radial-gradient(circle at 80% 20%, rgba(255, 255, 255, 0.1) 0%, transparent 50%);
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
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  padding: 20px 30px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.title-section {
  flex: 1;
  text-align: center;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 800;
  background: linear-gradient(135deg, #fff 0%, #e2e8f0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0 0 8px 0;
}

.subtitle {
  color: rgba(255, 255, 255, 0.8);
  font-size: 1.1rem;
  margin: 0;
}

/* Состояния */
.loading-container,
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 60px 40px;
  text-align: center;
}

.spinner {
  width: 60px;
  height: 60px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-top-color: #4f46e5;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 24px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.retry-btn {
  margin-top: 20px;
  padding: 12px 32px;
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.retry-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(79, 70, 229, 0.4);
}

/* Leaderboard */
.leaderboard-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.15);
}

.leaderboard-table {
  max-height: 70vh;
  overflow-y: auto;
}

.table-header {
  display: grid;
  grid-template-columns: 80px 1fr 200px;
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  color: white;
  padding: 20px 32px;
  font-weight: 700;
  font-size: 1.1rem;
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-cell {
  font-weight: 700;
  display: flex;
  align-items: center;
}

.table-row {
  display: grid;
  grid-template-columns: 80px 1fr 200px;
  padding: 20px 32px;
  border-bottom: 1px solid #f1f5f9;
  transition: all 0.3s ease;
  cursor: pointer;
}

.table-row:hover {
  background: #f8fafc;
  transform: translateX(8px);
}

.table-row.top-3 {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  font-weight: 600;
}

.table-row.top-3:hover {
  background: linear-gradient(135deg, #fed7aa 0%, #fdc08a 100%);
}

.table-row.current-user {
  border-left: 4px solid #4f46e5;
  background: linear-gradient(90deg, rgba(79, 70, 229, 0.1) 0%, transparent 20px);
}

.cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.rank {
  font-size: 1.3rem;
  font-weight: 700;
  color: #64748b;
  justify-content: center;
}

.medal {
  font-size: 1.8rem;
}

.username {
  font-weight: 600;
  font-size: 1.1rem;
  color: #1e293b;
}

.current-badge {
  background: #4f46e5;
  color: white;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
  margin-left: 12px;
}

.score {
  font-weight: 700;
  font-size: 1.3rem;
  color: #4f46e5;
  justify-content: flex-end;
}

/* Пустое состояние */
.empty-state {
  padding: 80px 40px;
  text-align: center;
  color: #64748b;
}

.empty-icon {
  font-size: 6rem;
  display: block;
  margin-bottom: 24px;
  opacity: 0.5;
}

/* Информация */
.info {
  padding: 24px 32px;
  background: #f8fafc;
  text-align: center;
  color: #64748b;
  font-size: 0.95rem;
}

/* Адаптивность */
@media (max-width: 768px) {
  .content-wrapper {
    padding: 16px;
  }
  
  .header {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
  
  .page-title {
    font-size: 2rem;
  }
  
  .table-header,
  .table-row {
    grid-template-columns: 60px 1fr 120px;
    padding: 16px 20px;
    font-size: 0.95rem;
  }
  
  .medal {
    font-size: 1.4rem;
  }
}

@media (max-width: 480px) {
  .table-header,
  .table-row {
    grid-template-columns: 50px 1fr;
    gap: 16px;
  }
  
  .score {
    grid-column: 1 / -1;
    justify-content: flex-start;
    margin-top: 8px;
    font-size: 1.1rem;
  }
}
</style>
