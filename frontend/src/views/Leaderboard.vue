<template>
  <div class="leaderboard-page">
    <!-- Фон -->
    <div class="background"></div>
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
/* Основные стили страницы */
.leaderboard-page {
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
  max-width: 1000px;
  margin: 0 auto;
  text-align: center;
}

/* Шапка страницы - ИСПРАВЛЕНА! */
.header {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 40px;
  padding: 25px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.2);
  position: relative;
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
  flex-shrink: 0;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.btn-icon {
  font-size: 1.2rem;
}

/* Центрированный заголовок */
.title-section {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  text-align: center;
  pointer-events: none;
}

.page-title {
  font-size: 2.2rem;
  color: white;
  margin: 0 0 10px 0;
  font-weight: 700;
  text-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
  white-space: nowrap;
}

.subtitle {
  font-size: 1.1rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  font-weight: 400;
}

/* Невидимая кнопка-плейсхолдер для симметрии */
.placeholder-btn {
  width: 120px;
  height: 48px;
  visibility: hidden;
  flex-shrink: 0;
}

/* Состояние загрузки */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 60px 40px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.2);
  width: 100%;
  margin: 20px 0;
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

.loading-container p {
  color: rgba(255, 255, 255, 0.9);
  font-size: 1.1rem;
  margin: 0;
}

/* Состояние ошибки */
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  padding: 50px 40px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(239, 68, 68, 0.3);
  width: 100%;
  max-width: 600px;
  margin: 20px 0;
}

.error-icon {
  font-size: 3.5rem;
  margin-bottom: 10px;
}

.error-container h3 {
  color: white;
  font-size: 1.5rem;
  margin: 0 0 10px 0;
}

.error-container p {
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

/* Контейнер лидерборда - БЕЗ ПОЛОСЫ ПРОКРУТКИ */
.leaderboard-container {
  width: 100%;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  backdrop-filter: blur(10px);
  border: 2px solid rgba(255, 255, 255, 0.2);
  padding: 30px;
  margin-bottom: 30px;
  overflow: hidden; /* Важно: скрываем переполнение */
}

/* Таблица лидеров - ИСПРАВЛЕНО: БЕЗ ПОЛОСЫ ПРОКРУТКИ */
.leaderboard-table {
  max-height: 600px;
  overflow-y: auto; /* Только вертикальная прокрутка */
  border-radius: 15px;
  margin-bottom: 20px;
  
  /* Скрываем полосу прокрутки во всех браузерах */
  scrollbar-width: none; /* Firefox */
  -ms-overflow-style: none; /* IE и Edge */
}

/* Скрываем полосу прокрутки в WebKit браузерах (Chrome, Safari) */
.leaderboard-table::-webkit-scrollbar {
  display: none;
  width: 0;
  height: 0;
}

.table-header {
  display: grid;
  grid-template-columns: 80px 1fr 150px;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  padding: 18px 20px;
  font-weight: 600;
  font-size: 1.1rem;
  border-radius: 15px 15px 0 0;
  position: sticky;
  top: 0;
  z-index: 10;
  backdrop-filter: blur(10px);
}

.header-cell {
  display: flex;
  align-items: center;
}

.rank {
  justify-content: center;
}

.score {
  justify-content: flex-end;
}

/* Строки таблицы - УЛУЧШЕННЫЙ HOVER ЭФФЕКТ */
.table-row {
  display: grid;
  grid-template-columns: 80px 1fr 150px;
  padding: 18px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.05);
  position: relative;
}

/* Hover эффект - ПЛАВНОЕ ВЫДЕЛЕНИЕ */
.table-row:hover {
  background: rgba(255, 255, 255, 0.2);
  box-shadow: 
    0 4px 12px rgba(0, 0, 0, 0.1),
    0 0 0 2px rgba(255, 255, 255, 0.3);
  transform: scale(1.01);
  z-index: 1;
}

/* Анимация для плавного выделения */
.table-row::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(255, 255, 255, 0.1) 50%,
    transparent 100%
  );
  opacity: 0;
  transition: opacity 0.3s ease;
  pointer-events: none;
  border-radius: 8px;
}

.table-row:hover::before {
  opacity: 1;
}

/* Топ-3 места с улучшенным hover */
.table-row.top-3 {
  background: linear-gradient(90deg, 
    rgba(255, 215, 0, 0.15) 0%, 
    rgba(192, 192, 192, 0.1) 33%, 
    rgba(205, 127, 50, 0.1) 66%, 
    transparent 100%);
}

.table-row.top-3:hover {
  background: linear-gradient(90deg, 
    rgba(255, 215, 0, 0.25) 0%, 
    rgba(192, 192, 192, 0.2) 33%, 
    rgba(205, 127, 50, 0.2) 66%, 
    rgba(255, 255, 255, 0.2) 100%);
  box-shadow: 
    0 4px 15px rgba(255, 215, 0, 0.3),
    0 0 0 2px rgba(255, 215, 0, 0.5);
}

.table-row.top-3:nth-child(1) {
  background: linear-gradient(90deg, 
    rgba(255, 215, 0, 0.25) 0%, 
    rgba(255, 215, 0, 0.15) 50%, 
    transparent 100%);
}

.table-row.top-3:nth-child(1):hover {
  background: linear-gradient(90deg, 
    rgba(255, 215, 0, 0.35) 0%, 
    rgba(255, 215, 0, 0.25) 50%, 
    rgba(255, 255, 255, 0.2) 100%);
}

.table-row.top-3:nth-child(2) {
  background: linear-gradient(90deg, 
    rgba(192, 192, 192, 0.2) 0%, 
    rgba(192, 192, 192, 0.1) 50%, 
    transparent 100%);
}

.table-row.top-3:nth-child(2):hover {
  background: linear-gradient(90deg, 
    rgba(192, 192, 192, 0.3) 0%, 
    rgba(192, 192, 192, 0.2) 50%, 
    rgba(255, 255, 255, 0.2) 100%);
}

.table-row.top-3:nth-child(3) {
  background: linear-gradient(90deg, 
    rgba(205, 127, 50, 0.18) 0%, 
    rgba(205, 127, 50, 0.08) 50%, 
    transparent 100%);
}

.table-row.top-3:nth-child(3):hover {
  background: linear-gradient(90deg, 
    rgba(205, 127, 50, 0.28) 0%, 
    rgba(205, 127, 50, 0.18) 50%, 
    rgba(255, 255, 255, 0.2) 100%);
}

/* Текущий пользователь */
.table-row.current-user {
  background: linear-gradient(90deg, 
    rgba(79, 70, 229, 0.25) 0%, 
    rgba(79, 70, 229, 0.15) 50%, 
    transparent 100%);
  border-left: 4px solid #4f46e5;
}

.table-row.current-user:hover {
  background: linear-gradient(90deg, 
    rgba(79, 70, 229, 0.35) 0%, 
    rgba(79, 70, 229, 0.25) 50%, 
    rgba(255, 255, 255, 0.2) 100%);
  box-shadow: 
    0 4px 15px rgba(79, 70, 229, 0.4),
    0 0 0 2px rgba(79, 70, 229, 0.6);
}

.cell {
  display: flex;
  align-items: center;
  color: white;
  position: relative;
  z-index: 2; /* Чтобы текст был над hover эффектом */
}

.cell.rank {
  font-size: 1.3rem;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
  justify-content: center;
}

.medal {
  font-size: 1.8rem;
}

.cell.name {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username {
  font-weight: 600;
  font-size: 1.1rem;
  transition: transform 0.3s ease;
}

.table-row:hover .username {
  transform: translateX(5px);
}

.current-badge {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: 600;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.table-row:hover .current-badge {
  transform: scale(1.1);
  box-shadow: 0 2px 8px rgba(79, 70, 229, 0.5);
}

.cell.score {
  font-weight: 700;
  font-size: 1.2rem;
  color: #fbbf24;
  justify-content: flex-end;
  transition: transform 0.3s ease;
}

.table-row:hover .cell.score {
  transform: scale(1.1);
  color: #ffd700;
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
  margin: 30px 0;
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
  margin: 0;
  font-size: 1.1rem;
  max-width: 400px;
}

/* Информация о времени обновления */
.info {
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.95rem;
  text-align: center;
  padding: 15px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

/* Адаптивность */
@media (max-width: 768px) {
  .content-wrapper {
    padding: 20px 15px;
  }
  
  .header {
    flex-direction: row;
    justify-content: space-between;
    gap: 10px;
    text-align: center;
    padding: 20px;
  }
  
  .title-section {
    position: static;
    transform: none;
    order: 2;
    margin: 0 10px;
    flex-grow: 1;
  }
  
  .back-btn {
    order: 1;
    padding: 10px 20px;
  }
  
  .page-title {
    font-size: 1.8rem;
    white-space: normal;
  }
  
  .leaderboard-container {
    padding: 20px;
  }
  
  .table-header,
  .table-row {
    grid-template-columns: 60px 1fr 120px;
    padding: 15px;
    font-size: 0.95rem;
  }
  
  .table-row:hover {
    transform: scale(1.02);
  }
  
  .medal {
    font-size: 1.4rem;
  }
  
  .username {
    font-size: 1rem;
  }
  
  .cell.score {
    font-size: 1.1rem;
  }
}

@media (max-width: 480px) {
  .content-wrapper {
    padding: 15px 10px;
  }
  
  .header {
    padding: 15px;
    flex-wrap: wrap;
  }
  
  .back-btn {
    width: 100%;
    justify-content: center;
    order: 1;
    margin-bottom: 10px;
  }
  
  .title-section {
    order: 2;
    width: 100%;
    margin: 0;
  }
  
  .page-title {
    font-size: 1.6rem;
    margin-bottom: 5px;
  }
  
  .subtitle {
    font-size: 1rem;
  }
  
  .table-header,
  .table-row {
    grid-template-columns: 50px 1fr 100px;
    padding: 12px 10px;
    font-size: 0.9rem;
  }
  
  .empty-state {
    padding: 50px 20px;
  }
  
  .empty-icon {
    font-size: 4rem;
  }
}
</style>
