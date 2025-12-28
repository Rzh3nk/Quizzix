<template>
  <div class="quiz-detail-page">
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
        <router-link 
          v-if="quiz?.category" 
          :to="`/categories/${quiz.category_id}`" 
          class="breadcrumb-link"
        >
          {{ quiz.category.name }}
        </router-link>
        <span v-else class="breadcrumb-current">Загрузка...</span>
        <span class="breadcrumb-separator">›</span>
        <span class="breadcrumb-current">{{ quiz?.title || 'Квиз' }}</span>
      </div>

      <!-- Основной контент -->
      <div class="quiz-detail-container">
        <!-- Состояние загрузки -->
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>Загрузка квиза...</p>
        </div>

        <!-- Состояние ошибки -->
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
          <!-- Шапка квиза -->
          <div class="quiz-header">
            <!-- Изображение квиза -->
            <div class="quiz-image-container">
              <div v-if="quiz.ImgURL" class="quiz-image-wrapper">
                <img 
                  :src="getImageUrl(quiz.ImgURL)" 
                  :alt="quiz.title"
                  class="quiz-image"
                  @error="handleImageError"
                />
              </div>
              <div v-else class="quiz-image-placeholder">
                {{ getQuizIcon(quiz.title) }}
              </div>
              
              <!-- Бейдж сложности -->
              <div v-if="quiz.difficulty" class="quiz-difficulty" :class="quiz.difficulty">
                {{ getDifficultyText(quiz.difficulty) }}
              </div>
            </div>

            <!-- Основная информация -->
            <div class="quiz-info">
              <div class="quiz-meta">
                <span class="quiz-category">
                  <span class="category-icon">🏷️</span>
                  {{ quiz.category?.name || 'Без категории' }}
                </span>
                <span v-if="quiz.timeLimit" class="quiz-time">
                  <span class="time-icon">⏱️</span>
                  {{ quiz.timeLimit }} мин
                </span>
                <span class="quiz-created">
                  <span class="date-icon">📅</span>
                  {{ formatDate(quiz.createdAt) }}
                </span>
              </div>

              <h1 class="quiz-title">{{ quiz.title }}</h1>
              
              <div class="quiz-description">
                {{ quiz.description || 'Нет описания' }}
              </div>

              <!-- Статистика квиза -->
              <div class="quiz-stats">
                <div class="stat">
                  <div class="stat-icon">❓</div>
                  <div class="stat-content">
                    <div class="stat-value">{{ questionCount }}</div>
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
                
                <div v-if="quiz.rating" class="stat">
                  <div class="stat-icon">⭐</div>
                  <div class="stat-content">
                    <div class="stat-value">{{ quiz.rating.toFixed(1) }}</div>
                    <div class="stat-label">рейтинг</div>
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

          <!-- Детальная информация -->
          <div class="quiz-details">
            <!-- Сложность и требования -->
            <div class="details-section">
              <h3 class="section-title">
                <span class="section-icon">📊</span>
                Информация о квизе
              </h3>
              
              <div class="details-grid">
                <div class="detail-item">
                  <span class="detail-icon">⚡</span>
                  <div class="detail-content">
                    <div class="detail-label">Сложность</div>
                    <div class="detail-value" :class="quiz.difficulty">
                      {{ getDifficultyText(quiz.difficulty || 'medium') }}
                    </div>
                  </div>
                </div>
                
                <div class="detail-item">
                  <span class="detail-icon">⏱️</span>
                  <div class="detail-content">
                    <div class="detail-label">Время</div>
                    <div class="detail-value">
                      {{ quiz.timeLimit ? `${quiz.timeLimit} минут` : 'Без ограничения' }}
                    </div>
                  </div>
                </div>
                
                <div class="detail-item">
                  <span class="detail-icon">❓</span>
                  <div class="detail-content">
                    <div class="detail-label">Вопросы</div>
                    <div class="detail-value">{{ questionCount }}</div>
                  </div>
                </div>
                
                <div class="detail-item">
                  <span class="detail-icon">🎯</span>
                  <div class="detail-content">
                    <div class="detail-label">Тип вопросов</div>
                    <div class="detail-value">С выбором ответа</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Советы по прохождению -->
            <div class="tips-section">
              <h3 class="section-title">
                <span class="section-icon">💡</span>
                Советы по прохождению
              </h3>
              
              <div class="tips-list">
                <div class="tip-item">
                  <span class="tip-icon">👀</span>
                  <div class="tip-content">Внимательно читайте вопросы и варианты ответов</div>
                </div>
                
                <div class="tip-item">
                  <span class="tip-icon">⏱️</span>
                  <div class="tip-content">Следите за временем, если установлен лимит</div>
                </div>
                
                <div class="tip-item">
                  <span class="tip-icon">🤔</span>
                  <div class="tip-content">Не спешите отвечать - подумайте над каждым вопросом</div>
                </div>
              </div>
            </div>

            <!-- Лучшие результаты -->
            <div v-if="topResults.length > 0" class="leaderboard-section">
              <h3 class="section-title">
                <span class="section-icon">🏆</span>
                Лучшие результаты
              </h3>
              
              <div class="leaderboard">
                <div 
                  v-for="(result, index) in topResults" 
                  :key="result.id" 
                  class="leaderboard-item"
                  :class="{ 'first': index === 0, 'second': index === 1, 'third': index === 2 }"
                >
                  <div class="rank">
                    <span class="rank-number">{{ index + 1 }}</span>
                    <span class="rank-icon">{{ getRankIcon(index) }}</span>
                  </div>
                  
                  <div class="player-info">
                    <div class="player-name">{{ result.userName || 'Аноним' }}</div>
                    <div class="player-score">{{ result.score }}%</div>
                  </div>
                  
                  <div class="result-info">
                    <div class="result-time">{{ formatDuration(result.timeSpent) }}</div>
                    <div class="result-date">{{ formatDate(result.createdAt) }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Кнопки действий -->
          <div class="action-buttons">
            <button @click="startQuiz" class="start-btn">
              <span class="btn-icon">🚀</span>
              <span class="btn-text">Начать квиз</span>
            </button>
            
            <button @click="goBack" class="back-btn">
              <span class="btn-icon">←</span>
              <span class="btn-text">Назад</span>
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

const route = useRoute()
const router = useRouter()

//const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Реактивные данные
const quiz = ref(null)
const loading = ref(true)
const error = ref(null)
const topResults = ref([])

// Компьютеды
const quizId = computed(() => route.params.id)

const questionCount = computed(() => {
  return quiz.value?.Questions?.length || 0
})

// Методы
const fetchQuiz = async () => {
  try {
    loading.value = true
    error.value = null
    
    console.log(quizId.value)
    const response = await fetch(`/api/quiz/${quizId.value}`)
    
    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('Квиз не найден')
      }
      throw new Error('Ошибка загрузки квиза')
    }
    
    quiz.value = await response.json()
    
    // Загружаем топ результаты
    await fetchTopResults()
    
  } catch (err) {
    console.error('Ошибка загрузки квиза:', err)
    error.value = err.message || 'Не удалось загрузить квиз'
    
    // Мок данные для разработки
    if (import.meta.env.DEV) {
      quiz.value = createMockQuiz()
      topResults.value = createMockResults()
    }
  } finally {
    loading.value = false
  }
}

const fetchTopResults = async () => {
  try {
    const response = await fetch(`/api/quiz/${quizId.value}/results/top`)
    if (response.ok) {
      topResults.value = await response.json()
    }
  } catch (err) {
    console.error('Ошибка загрузки результатов:', err)
    // Используем мок данные
    topResults.value = createMockResults()
  }
}

const createMockQuiz = () => {
  return {
    id: quizId.value,
    title: 'История Древнего Рима',
    description: 'Проверьте свои знания о Римской империи, её императорах, завоеваниях и культуре. От основания города до падения империи.',
    ImgURL: null,
    difficulty: 'medium',
    timeLimit: 20,
    plays: 1245,
    rating: 4.7,
    authorName: 'Иван Иванов',
    authorStats: {
      quizCount: 15,
      totalPlays: 8900
    },
    category: {
      id: 1,
      name: 'История'
    },
    category_id: 1,
    createdAt: new Date().toISOString(),
    Questions: Array.from({ length: 15 }, (_, i) => ({
      id: i + 1,
      text: `Вопрос ${i + 1}`,
      answers: []
    }))
  }
}

const createMockResults = () => {
  return [
    { id: 1, userName: 'Алексей Петров', score: 98, timeSpent: 850, createdAt: new Date(Date.now() - 86400000).toISOString() },
    { id: 2, userName: 'Мария Смирнова', score: 95, timeSpent: 920, createdAt: new Date(Date.now() - 172800000).toISOString() },
    { id: 3, userName: 'Дмитрий Иванов', score: 93, timeSpent: 780, createdAt: new Date(Date.now() - 259200000).toISOString() },
    { id: 4, userName: 'Анна Козлова', score: 90, timeSpent: 1100, createdAt: new Date(Date.now() - 345600000).toISOString() },
    { id: 5, userName: 'Сергей Волков', score: 88, timeSpent: 950, createdAt: new Date(Date.now() - 432000000).toISOString() }
  ]
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

const getDifficultyText = (difficulty) => {
  const map = { 
    easy: 'Лёгкий', 
    medium: 'Средний', 
    hard: 'Сложный' 
  }
  return map[difficulty] || 'Средний'
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

const formatNumber = (num) => {
  return new Intl.NumberFormat('ru-RU').format(num)
}

const formatDuration = (seconds) => {
  if (!seconds) return '0 сек'
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return mins > 0 ? `${mins} мин ${secs} сек` : `${secs} сек`
}

const getRankIcon = (index) => {
  return ['🥇', '🥈', '🥉'][index] || '🏅'
}

const startQuiz = () => {
  router.push(`/quiz/${quizId.value}/play`)
}

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/categories')
  }
}

// Наблюдатель за изменением ID квиза
import { watch } from 'vue'

watch(quizId, (newId) => {
  if (newId) {
    fetchQuiz()
  }
})

// Хуки жизненного цикла
onMounted(() => {
  if (quizId.value) {
    fetchQuiz()
  }
})
</script>

<style scoped>
/* Основные стили страницы */
.quiz-detail-page {
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
  background: #f8fafc;
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
}

/* Хлебные крошки */
.breadcrumbs {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 30px;
  width: 100%;
  color: white;
  flex-wrap: wrap;
}

.breadcrumb-link {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: white;
  text-decoration: none;
  opacity: 0.8;
  transition: opacity 0.3s ease;
  padding: 8px 12px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
}

.breadcrumb-link:hover {
  opacity: 1;
  background: rgba(255, 255, 255, 0.2);
}

.breadcrumb-separator {
  color: rgba(255, 255, 255, 0.6);
  margin: 0 5px;
}

.breadcrumb-current {
  color: white;
  font-weight: 500;
  padding: 8px 12px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.15);
}

.breadcrumb-icon {
  font-size: 1.1em;
}

/* Контейнер квиза */
.quiz-detail-container {
  width: 100%;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

/* Состояния загрузки и ошибки */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
}

.spinner {
  width: 60px;
  height: 60px;
  border: 4px solid #e5e7eb;
  border-top-color: #4f46e5;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-state p {
  color: #6b7280;
  font-size: 1.1rem;
}

.error-state {
  text-align: center;
  padding: 60px 20px;
}

.error-icon {
  font-size: 4rem;
  color: #ef4444;
  margin-bottom: 20px;
  display: block;
}

.error-state h3 {
  font-size: 1.8rem;
  color: #1f2937;
  margin-bottom: 10px;
}

.error-state p {
  color: #6b7280;
  margin-bottom: 30px;
  font-size: 1.1rem;
}

/* Шапка квиза */
.quiz-header {
  display: flex;
  gap: 40px;
  margin-bottom: 40px;
  flex-wrap: wrap;
}

.quiz-image-container {
  flex: 0 0 300px;
  position: relative;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
}

.quiz-image-wrapper {
  width: 100%;
  height: 200px;
  overflow: hidden;
}

.quiz-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.quiz-image:hover {
  transform: scale(1.05);
}

.quiz-image-placeholder {
  width: 100%;
  height: 200px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 4rem;
  color: white;
}

.quiz-difficulty {
  position: absolute;
  top: 15px;
  right: 15px;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 600;
  text-transform: uppercase;
  color: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.quiz-difficulty.easy {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
}

.quiz-difficulty.medium {
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
}

.quiz-difficulty.hard {
  background: linear-gradient(135deg, #ef4444 0%, #f87171 100%);
}

/* Основная информация о квизе */
.quiz-info {
  flex: 1;
  min-width: 300px;
}

.quiz-meta {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

.quiz-category,
.quiz-time,
.quiz-created {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #f3f4f6;
  border-radius: 10px;
  color: #4b5563;
  font-size: 0.9rem;
  font-weight: 500;
}

.quiz-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 20px;
  color: #1f2937;
  line-height: 1.2;
}

.quiz-description {
  font-size: 1.1rem;
  line-height: 1.6;
  color: #4b5563;
  margin-bottom: 30px;
  padding: 20px;
  background: #f9fafb;
  border-radius: 12px;
  border-left: 4px solid #4f46e5;
}

/* Статистика квиза */
.quiz-stats {
  display: flex;
  gap: 30px;
  margin-bottom: 30px;
  flex-wrap: wrap;
}

.stat {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  border: 2px solid #e5e7eb;
  min-width: 160px;
  transition: all 0.3s ease;
}

.stat:hover {
  border-color: #4f46e5;
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(79, 70, 229, 0.1);
}

.stat-icon {
  font-size: 2rem;
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #eef2ff;
  border-radius: 10px;
  color: #4f46e5;
}

.stat-content {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 1.8rem;
  font-weight: 700;
  color: #1f2937;
  line-height: 1;
}

.stat-label {
  font-size: 0.9rem;
  color: #6b7280;
  margin-top: 5px;
}

/* Автор квиза */
.quiz-author {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 20px;
  background: white;
  border-radius: 12px;
  border: 2px solid #e5e7eb;
}

.author-avatar {
  width: 50px;
  height: 50px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.5rem;
}

.author-info {
  flex: 1;
}

.author-name {
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 5px;
}

.author-stats {
  font-size: 0.9rem;
  color: #6b7280;
}

/* Детальная информация */
.quiz-details {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
  margin-bottom: 40px;
}

.details-section,
.tips-section,
.leaderboard-section {
  background: white;
  border-radius: 16px;
  padding: 30px;
  border: 2px solid #e5e7eb;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1.3rem;
  margin-bottom: 25px;
  color: #1f2937;
  padding-bottom: 15px;
  border-bottom: 2px solid #f3f4f6;
}

.section-icon {
  font-size: 1.3em;
}

/* Детали квиза */
.details-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  background: #f9fafb;
  border-radius: 10px;
  transition: all 0.3s ease;
}

.detail-item:hover {
  background: #f3f4f6;
  transform: translateX(5px);
}

.detail-icon {
  font-size: 1.5rem;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #eef2ff;
  border-radius: 8px;
  color: #4f46e5;
}

.detail-content {
  flex: 1;
}

.detail-label {
  font-size: 0.9rem;
  color: #6b7280;
  margin-bottom: 4px;
}

.detail-value {
  font-weight: 600;
  color: #1f2937;
}

.detail-value.easy {
  color: #10b981;
}

.detail-value.medium {
  color: #f59e0b;
}

.detail-value.hard {
  color: #ef4444;
}

/* Советы */
.tips-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.tip-item {
  display: flex;
  align-items: flex-start;
  gap: 15px;
  padding: 15px;
  background: #f0f9ff;
  border-radius: 10px;
  border-left: 4px solid #0ea5e9;
}

.tip-icon {
  font-size: 1.2rem;
  margin-top: 2px;
  color: #0ea5e9;
}

.tip-content {
  color: #1e40af;
  font-size: 0.95rem;
  line-height: 1.5;
}

/* Таблица лидеров */
.leaderboard {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.leaderboard-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  background: #f9fafb;
  border-radius: 10px;
  transition: all 0.3s ease;
}

.leaderboard-item:hover {
  background: #f3f4f6;
  transform: translateX(5px);
}

.leaderboard-item.first {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  border: 2px solid #fbbf24;
}

.leaderboard-item.second {
  background: linear-gradient(135deg, #f3f4f6 0%, #e5e7eb 100%);
  border: 2px solid #d1d5db;
}

.leaderboard-item.third {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  border: 2px solid #fbbf24;
  opacity: 0.9;
}

.rank {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 40px;
}

.rank-number {
  font-size: 1.2rem;
  font-weight: 700;
  color: #4f46e5;
}

.rank-icon {
  font-size: 1rem;
}

.player-info {
  flex: 1;
}

.player-name {
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.player-score {
  font-size: 1.1rem;
  font-weight: 700;
  color: #10b981;
}

.result-info {
  text-align: right;
  min-width: 120px;
}

.result-time {
  font-size: 0.9rem;
  color: #6b7280;
  margin-bottom: 4px;
}

.result-date {
  font-size: 0.8rem;
  color: #9ca3af;
}

/* Кнопки действий */
.action-buttons {
  display: flex;
  gap: 20px;
  justify-content: center;
  padding-top: 40px;
  border-top: 2px solid #f3f4f6;
}

.start-btn,
.back-btn {
  padding: 18px 40px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 1.1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 12px;
  border: none;
  min-width: 200px;
  justify-content: center;
}

.start-btn {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  box-shadow: 0 10px 25px rgba(79, 70, 229, 0.3);
}

.start-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 15px 35px rgba(79, 70, 229, 0.4);
}

.back-btn {
  background: white;
  color: #4f46e5;
  border: 2px solid #4f46e5;
}

.back-btn:hover {
  background: #f8fafc;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(79, 70, 229, 0.1);
}

.btn-icon {
  font-size: 1.3em;
}

.btn-text {
  font-size: 1.1rem;
}

/* Адаптивность */
@media (max-width: 992px) {
  .quiz-header {
    flex-direction: column;
  }
  
  .quiz-image-container {
    flex: none;
    width: 100%;
    max-width: 400px;
    margin: 0 auto;
  }
  
  .details-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .content-wrapper {
    padding: 20px 15px;
  }
  
  .quiz-detail-container {
    padding: 25px;
  }
  
  .quiz-title {
    font-size: 2rem;
  }
  
  .quiz-stats {
    justify-content: center;
  }
  
  .stat {
    min-width: 140px;
  }
  
  .action-buttons {
    flex-direction: column;
    align-items: center;
  }
  
  .start-btn,
  .back-btn {
    width: 100%;
    max-width: 300px;
  }
  
  .breadcrumbs {
    font-size: 0.9rem;
  }
  
  .breadcrumb-link,
  .breadcrumb-current {
    padding: 6px 10px;
  }
}

@media (max-width: 480px) {
  .quiz-detail-container {
    padding: 20px;
  }
  
  .quiz-title {
    font-size: 1.8rem;
  }
  
  .quiz-stats {
    flex-direction: column;
    align-items: center;
  }
  
  .stat {
    width: 100%;
    max-width: 250px;
  }
  
  .section-title {
    font-size: 1.2rem;
  }
  
  .details-section,
  .tips-section,
  .leaderboard-section {
    padding: 20px;
  }
}
</style>