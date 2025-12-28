<template>
  <div class="quiz-play-page">
    <!-- Фон -->
    <div class="background"></div>

    <!-- Контент -->
    <div class="content-wrapper">
      <!-- Шапка с прогрессом -->
      <div class="quiz-header">
        <button @click="exitQuiz" class="exit-btn">
          <span class="btn-icon">←</span>
          Выйти
        </button>
        
        <div class="progress-info">
          <h1 class="quiz-title">{{ quiz?.Title || 'Квиз' }}</h1>
          <div class="progress">
            <span class="progress-text">Вопрос {{ currentQuestionNumber }} из {{ questions.length }}</span>
            <div class="progress-bar">
              <div 
                class="progress-fill" 
                :style="{ width: progressPercentage + '%' }"
              ></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Основной контейнер -->
      <div class="main-container">
        <!-- Загрузка -->
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>Загрузка вопросов...</p>
        </div>

        <!-- Ошибка -->
        <div v-else-if="error" class="error-state">
          <span class="error-icon">⚠️</span>
          <h3>Ошибка загрузки</h3>
          <p>{{ error }}</p>
          <button @click="goBack" class="back-btn">
            Вернуться к квизу
          </button>
        </div>

        <!-- Результаты -->
        <div v-else-if="showResults" class="results-container">
          <div class="results-card">
            <div class="results-header">
              <span class="results-icon">🏆</span>
              <h2 class="results-title">Квиз завершен!</h2>
            </div>

            <div class="score-display">
              <div class="score-circle">
                <div class="score-value">{{ resultData.percent }}%</div>
                <div class="score-text">точность</div>
              </div>

              <div class="score-details">
                <div class="score-stat">
                  <div class="stat-value">{{ resultData.score }}</div>
                  <div class="stat-label">правильных ответов</div>
                </div>
                <div class="score-stat">
                  <div class="stat-value">{{ resultData.total }}</div>
                  <div class="stat-label">всего вопросов</div>
                </div>
              </div>
            </div>

            <!-- Детали по вопросам -->
            <div class="questions-results">
              <h3 class="section-title">Детали ответов</h3>
              
              <div class="questions-list">
                <div 
                  v-for="(detail, index) in resultData.details" 
                  :key="detail.question_id"
                  class="question-result"
                  :class="{ 'correct': detail.is_correct, 'incorrect': !detail.is_correct }"
                >
                  <div class="question-number">Вопрос {{ index + 1 }}</div>
                  <div class="question-status">
                    <span v-if="detail.is_correct" class="status-icon correct">✓</span>
                    <span v-else class="status-icon incorrect">✗</span>
                    {{ detail.is_correct ? 'Правильно' : 'Неправильно' }}
                  </div>
                </div>
              </div>
            </div>

            <div class="results-actions">
              <button @click="restartQuiz" class="action-btn primary">
                <span class="btn-icon">🔄</span>
                Пройти еще раз
              </button>
              
              <button @click="goToQuizPage" class="action-btn">
                <span class="btn-icon">📚</span>
                К странице квиза
              </button>
            </div>
          </div>
        </div>

        <!-- Вопрос -->
        <div v-else-if="currentQuestion" class="question-container">
          <div class="question-card">
            <!-- Заголовок вопроса -->
            <div class="question-header">
              <div class="question-meta">
                <span class="question-number">Вопрос {{ currentQuestionNumber }}</span>
                <span v-if="currentQuestion.points" class="question-points">
                  <span class="points-icon">⭐</span>
                  {{ currentQuestion.points }} баллов
                </span>
              </div>
            </div>

            <!-- Текст вопроса -->
            <div class="question-content">
              <h2 class="question-text">{{ currentQuestion.text || currentQuestion.Text }}</h2>
              
              <div v-if="currentQuestion.image" class="question-image">
                <img :src="getImageUrl(currentQuestion.image)" :alt="currentQuestion.text">
              </div>
            </div>

            <!-- Варианты ответов -->
            <div class="answers-section">
              <h3 class="answers-title">Выберите правильный ответ:</h3>
              
              <div class="answers-list">
                <div 
                  v-for="answer in currentQuestion.answers" 
                  :key="answer.id || answer.ID"
                  class="answer-option"
                  :class="{ 
                    'selected': isAnswerSelected(answer),
                    'multiple': currentQuestion.isMultiple || false
                  }"
                  @click="toggleAnswer(answer)"
                >
                  <div class="answer-selector">
                    <div v-if="currentQuestion.isMultiple" class="checkbox">
                      <div class="checkbox-inner" :class="{ 'checked': isAnswerSelected(answer) }"></div>
                    </div>
                    <div v-else class="radio">
                      <div class="radio-inner" :class="{ 'checked': isAnswerSelected(answer) }"></div>
                    </div>
                  </div>
                  
                  <div class="answer-content">
                    <div class="answer-text">{{ answer.text || answer.Text }}</div>
                    
                    <div v-if="answer.image" class="answer-image">
                      <img :src="getImageUrl(answer.image)" :alt="answer.text">
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Кнопки навигации -->
            <div class="navigation-buttons">
              <button 
                v-if="!isFirstQuestion" 
                @click="prevQuestion" 
                class="nav-btn prev"
              >
                <span class="btn-icon">←</span>
                Назад
              </button>
              
              <div class="main-action">
                <button 
                  v-if="!isLastQuestion" 
                  @click="nextQuestion" 
                  class="next-btn"
                  :disabled="!hasSelectedAnswers"
                >
                  Следующий вопрос
                  <span class="btn-icon">→</span>
                </button>
                
                <button 
                  v-else 
                  @click="submitQuiz" 
                  class="submit-btn"
                  :disabled="!hasSelectedAnswers || isSubmitting"
                >
                  <span v-if="isSubmitting" class="spinner small"></span>
                  <span v-else class="btn-icon">🏁</span>
                  {{ isSubmitting ? 'Отправка...' : 'Завершить квиз' }}
                </button>
              </div>
            </div>
          </div>

          <!-- Индикатор вопросов -->
          <div class="questions-indicator">
            <div 
              v-for="(question, index) in questions" 
              :key="question.id || question.ID"
              class="question-dot"
              :class="{
                'current': index === currentQuestionIndex,
                'answered': isQuestionAnswered(question),
                'pending': index > currentQuestionIndex
              }"
              @click="goToQuestion(index)"
              :title="`Вопрос ${index + 1}`"
            >
              {{ index + 1 }}
            </div>
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
const questions = ref([])
const loading = ref(true)
const error = ref(null)
const currentQuestionIndex = ref(0)
const userAnswers = ref({}) // { questionId: [answerIds] }
const isSubmitting = ref(false)
const showResults = ref(false)
const resultData = ref({
  score: 0,
  total: 0,
  percent: 0,
  details: []
})

// Получаем ID квиза
const quizId = computed(() => route.params.id)

// Компьютеды
const currentQuestion = computed(() => {
  return questions.value[currentQuestionIndex.value]
})

const currentQuestionNumber = computed(() => {
  return currentQuestionIndex.value + 1
})

const isFirstQuestion = computed(() => {
  return currentQuestionIndex.value === 0
})

const isLastQuestion = computed(() => {
  return currentQuestionIndex.value === questions.value.length - 1
})

const progressPercentage = computed(() => {
  if (questions.value.length === 0) return 0
  return ((currentQuestionIndex.value + 1) / questions.value.length) * 100
})

const hasSelectedAnswers = computed(() => {
  if (!currentQuestion.value) return false
  
  const questionId = currentQuestion.value.id || currentQuestion.value.ID
  return Array.isArray(userAnswers.value[questionId]) && userAnswers.value[questionId].length > 0
})

// Методы
const fetchQuiz = async () => {
  try {
    loading.value = true
    error.value = null
    
    if (!quizId.value) {
      throw new Error('ID квиза не указан')
    }
    
    // Загружаем квиз с вопросами
    let apiUrl
    if (import.meta.env.DEV && API_URL.includes('localhost:8080')) {
      apiUrl = `/api/quizzes/${quizId.value}`
    } else {
      apiUrl = `/api/quizzes/${quizId.value}`
    }
    
    const token = localStorage.getItem('token')
    const headers = {
      'Content-Type': 'application/json'
    }
    
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
    
    const response = await fetch(apiUrl, { headers })
    
    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('Квиз не найден')
      }
      throw new Error(`Ошибка ${response.status}`)
    }
    
    const data = await response.json()
    
    // Преобразуем данные
    quiz.value = {
      ID: data.id || data.ID,
      Title: data.title || data.Title,
      Description: data.description || data.Description
    }
    
    // Получаем вопросы (если они есть в ответе)
    questions.value = data.questions || data.Questions || []
    
    if (questions.value.length === 0) {
      throw new Error('В этом квизе нет вопросов')
    }
    
    // Инициализируем структуру ответов
    initializeAnswers()
    
  } catch (err) {
    console.error('Ошибка загрузки квиза:', err)
    error.value = err.message || 'Не удалось загрузить вопросы'
    
    // Мок данные для разработки
    if (import.meta.env.DEV) {
      quiz.value = {
        ID: quizId.value || 1,
        Title: 'История Древнего Рима',
        Description: 'Проверьте свои знания'
      }
      
      questions.value = createMockQuestions()
      initializeAnswers()
    }
  } finally {
    loading.value = false
  }
}

const createMockQuestions = () => {
  return [
    {
      id: 1,
      text: 'Кто был первым императором Рима?',
      points: 10,
      isMultiple: false,
      answers: [
        { id: 10, text: 'Юлий Цезарь' },
        { id: 11, text: 'Октавиан Август' },
        { id: 12, text: 'Нерон' },
        { id: 13, text: 'Константин Великий' }
      ]
    },
    {
      id: 2,
      text: 'Какие из этих событий произошли в Древнем Риме?',
      points: 15,
      isMultiple: true,
      answers: [
        { id: 14, text: 'Основание Рима' },
        { id: 15, text: 'Пунические войны' },
        { id: 16, text: 'Великое переселение народов' },
        { id: 17, text: 'Падение Западной Римской империи' }
      ]
    },
    {
      id: 3,
      text: 'Кто написал "Энеиду"?',
      points: 10,
      isMultiple: false,
      answers: [
        { id: 18, text: 'Вергилий' },
        { id: 19, text: 'Гомер' },
        { id: 20, text: 'Цицерон' },
        { id: 21, text: 'Овидий' }
      ]
    }
  ]
}

const initializeAnswers = () => {
  // Создаем пустую структуру для ответов
  const answers = {}
  questions.value.forEach(question => {
    const questionId = question.id || question.ID
    answers[questionId] = []
  })
  userAnswers.value = answers
}

const getImageUrl = (path) => {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `${API_URL}${path.startsWith('/') ? path : '/' + path}`
}

const isAnswerSelected = (answer) => {
  if (!currentQuestion.value) return false
  
  const questionId = currentQuestion.value.id || currentQuestion.value.ID
  const answerId = answer.id || answer.ID
  return userAnswers.value[questionId]?.includes(answerId) || false
}

const toggleAnswer = (answer) => {
  if (!currentQuestion.value) return
  
  const questionId = currentQuestion.value.id || currentQuestion.value.ID
  const answerId = answer.id || answer.ID
  const currentAnswers = [...(userAnswers.value[questionId] || [])]
  
  if (currentQuestion.value.isMultiple) {
    // Множественный выбор
    const index = currentAnswers.indexOf(answerId)
    if (index === -1) {
      currentAnswers.push(answerId)
    } else {
      currentAnswers.splice(index, 1)
    }
  } else {
    // Одиночный выбор
    currentAnswers.length = 0
    currentAnswers.push(answerId)
  }
  
  userAnswers.value = {
    ...userAnswers.value,
    [questionId]: currentAnswers
  }
}

const isQuestionAnswered = (question) => {
  const questionId = question.id || question.ID
  return userAnswers.value[questionId]?.length > 0 || false
}

const nextQuestion = () => {
  if (currentQuestionIndex.value < questions.value.length - 1) {
    currentQuestionIndex.value++
  }
}

const prevQuestion = () => {
  if (currentQuestionIndex.value > 0) {
    currentQuestionIndex.value--
  }
}

const goToQuestion = (index) => {
  if (index >= 0 && index < questions.value.length) {
    currentQuestionIndex.value = index
  }
}

const submitQuiz = async () => {
  try {
    isSubmitting.value = true
    
    // Получаем ID пользователя (если есть)
    const userId = getUserId()
    
    // Подготавливаем данные для отправки
    const payload = {
      user_id: userId,
      answers: userAnswers.value
    }
    
    console.log('Отправляемые ответы:', payload)
    
    // Отправляем ответы на сервер
    let apiUrl
    if (import.meta.env.DEV && API_URL.includes('localhost:8080')) {
      apiUrl = `/api/quizzes/${quizId.value}/submit`
    } else {
      apiUrl = `/api/quizzes/${quizId.value}/submit`
    }
    
    const token = localStorage.getItem('token')
    const headers = {
      'Content-Type': 'application/json'
    }
    
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
    
    const response = await fetch(apiUrl, {
      method: 'POST',
      headers,
      body: JSON.stringify(payload)
    })
    
    if (!response.ok) {
      throw new Error(`Ошибка ${response.status}`)
    }
    
    const data = await response.json()
    console.log('Результаты:', data)
    
    // Сохраняем результаты
    resultData.value = {
      score: data.score || 0,
      total: data.total || questions.value.length,
      percent: data.percent || 0,
      details: data.details || [],
      result_id: data.result_id
    }
    
    // Показываем результаты
    showResults.value = true
    
  } catch (err) {
    console.error('Ошибка при отправке ответов:', err)
    alert('Произошла ошибка при отправке ответов. Попробуйте еще раз.')
  } finally {
    isSubmitting.value = false
  }
}

const getUserId = () => {
  // Пытаемся получить ID пользователя из localStorage или другого места
  // Пока возвращаем 1 как в примере
  return 1
}

const restartQuiz = () => {
  // Сбрасываем все данные
  currentQuestionIndex.value = 0
  initializeAnswers()
  showResults.value = false
  resultData.value = {
    score: 0,
    total: 0,
    percent: 0,
    details: []
  }
}

const exitQuiz = () => {
  if (confirm('Вы уверены, что хотите выйти из квиза? Прогресс будет потерян.')) {
    goToQuizPage()
  }
}

const goToQuizPage = () => {
  router.push(`/quiz/${quizId.value}`)
}

const goBack = () => {
  router.back()
}

// Хуки жизненного цикла
onMounted(() => {
  if (quizId.value) {
    fetchQuiz()
  }
})
</script>

<style scoped>
/* Основные стили страницы */
.quiz-play-page {
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
  padding: 20px;
  max-width: 900px;
  margin: 0 auto;
}

/* Шапка */
.quiz-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
  padding: 15px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.exit-btn {
  padding: 10px 20px;
  background: #f3f4f6;
  color: #4b5563;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.exit-btn:hover {
  background: #e5e7eb;
  transform: translateY(-2px);
}

.progress-info {
  flex: 1;
}

.quiz-title {
  font-size: 1.3rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.progress {
  display: flex;
  align-items: center;
  gap: 15px;
}

.progress-text {
  font-size: 0.95rem;
  color: #4f46e5;
  font-weight: 500;
  white-space: nowrap;
}

.progress-bar {
  flex: 1;
  height: 6px;
  background: #e5e7eb;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  border-radius: 3px;
  transition: width 0.3s ease;
}

/* Основной контейнер */
.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* Состояния */
.loading-state,
.error-state,
.results-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
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

.spinner.small {
  width: 20px;
  height: 20px;
  border-width: 2px;
  margin: 0;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-state p {
  color: #6b7280;
  font-size: 1.1rem;
}

.error-icon {
  font-size: 4rem;
  color: #ef4444;
  margin-bottom: 20px;
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
  text-align: center;
}

.error-state .back-btn {
  padding: 12px 24px;
  background: #4f46e5;
  color: white;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.error-state .back-btn:hover {
  background: #4338ca;
  transform: translateY(-2px);
}

/* Результаты */
.results-card {
  width: 100%;
  max-width: 500px;
  text-align: center;
}

.results-header {
  margin-bottom: 30px;
}

.results-icon {
  font-size: 4rem;
  display: block;
  margin-bottom: 10px;
}

.results-title {
  font-size: 2.2rem;
  color: #1f2937;
  margin: 0;
}

.score-display {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 40px;
  margin-bottom: 40px;
  flex-wrap: wrap;
}

.score-circle {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 15px 35px rgba(79, 70, 229, 0.3);
}

.score-value {
  font-size: 2.5rem;
  font-weight: 700;
  line-height: 1;
}

.score-text {
  font-size: 1rem;
  opacity: 0.9;
}

.score-details {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.score-stat {
  text-align: center;
  padding: 15px 25px;
  background: white;
  border-radius: 12px;
  border: 2px solid #e5e7eb;
  min-width: 150px;
  transition: all 0.3s ease;
}

.score-stat:hover {
  border-color: #4f46e5;
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(79, 70, 229, 0.1);
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #4f46e5;
  margin-bottom: 5px;
}

.stat-label {
  color: #6b7280;
  font-size: 0.9rem;
}

/* Детали вопросов */
.questions-results {
  margin-bottom: 40px;
}

.section-title {
  font-size: 1.3rem;
  color: #1f2937;
  margin-bottom: 20px;
  text-align: left;
  padding-bottom: 10px;
  border-bottom: 2px solid #f3f4f6;
}

.questions-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.question-result {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  border-radius: 10px;
  transition: all 0.3s ease;
}

.question-result.correct {
  background: #d1fae5;
  border-left: 4px solid #10b981;
}

.question-result.incorrect {
  background: #fee2e2;
  border-left: 4px solid #ef4444;
}

.question-number {
  font-weight: 500;
  color: #1f2937;
}

.question-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
}

.status-icon {
  font-size: 1.2rem;
  font-weight: bold;
}

.status-icon.correct {
  color: #10b981;
}

.status-icon.incorrect {
  color: #ef4444;
}

/* Кнопки действий */
.results-actions {
  display: flex;
  gap: 15px;
  justify-content: center;
}

/* Контейнер вопроса */
.question-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Карточка вопроса */
.question-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 30px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  flex: 1;
  display: flex;
  flex-direction: column;
}

.question-header {
  margin-bottom: 25px;
}

.question-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 15px;
  border-bottom: 2px solid #f3f4f6;
}

.question-number {
  font-size: 1.1rem;
  font-weight: 600;
  color: #4f46e5;
  background: #eef2ff;
  padding: 8px 15px;
  border-radius: 20px;
}

.question-points {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 15px;
  background: #fef3c7;
  color: #92400e;
  border-radius: 20px;
  font-weight: 500;
  font-size: 0.9rem;
}

.points-icon {
  font-size: 1.1rem;
}

.question-content {
  margin-bottom: 30px;
}

.question-text {
  font-size: 1.4rem;
  line-height: 1.5;
  color: #1f2937;
  margin-bottom: 20px;
}

.question-image {
  max-width: 100%;
  border-radius: 10px;
  overflow: hidden;
  margin-top: 15px;
}

.question-image img {
  width: 100%;
  height: auto;
  display: block;
}

/* Варианты ответов */
.answers-section {
  flex: 1;
}

.answers-title {
  font-size: 1.1rem;
  color: #4b5563;
  margin-bottom: 20px;
  font-weight: 500;
}

.answers-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.answer-option {
  display: flex;
  align-items: flex-start;
  gap: 15px;
  padding: 18px 20px;
  background: #f9fafb;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.answer-option:hover {
  background: #f3f4f6;
  border-color: #d1d5db;
  transform: translateX(5px);
}

.answer-option.selected {
  background: #e0e7ff;
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.answer-selector {
  flex-shrink: 0;
  margin-top: 2px;
}

.checkbox,
.radio {
  width: 22px;
  height: 22px;
  border: 2px solid #d1d5db;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
}

.answer-option.selected .checkbox,
.answer-option.selected .radio {
  border-color: #4f46e5;
  background: #4f46e5;
}

.radio {
  border-radius: 50%;
}

.checkbox-inner,
.radio-inner {
  width: 12px;
  height: 12px;
  border-radius: 2px;
  background: white;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.radio-inner {
  border-radius: 50%;
}

.checkbox-inner.checked,
.radio-inner.checked {
  opacity: 1;
}

.answer-content {
  flex: 1;
}

.answer-text {
  font-size: 1.1rem;
  color: #1f2937;
  line-height: 1.5;
}

.answer-image {
  max-width: 200px;
  margin-top: 10px;
  border-radius: 8px;
  overflow: hidden;
}

.answer-image img {
  width: 100%;
  height: auto;
  display: block;
}

/* Кнопки навигации */
.navigation-buttons {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 2px solid #f3f4f6;
}

.nav-btn,
.next-btn,
.submit-btn,
.action-btn {
  padding: 14px 28px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 10px;
  border: none;
  min-width: 180px;
  justify-content: center;
}

.nav-btn.prev {
  background: white;
  color: #4f46e5;
  border: 2px solid #4f46e5;
}

.nav-btn.prev:hover {
  background: #f8fafc;
  transform: translateY(-2px);
}

.next-btn,
.submit-btn {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  box-shadow: 0 10px 25px rgba(79, 70, 229, 0.3);
}

.next-btn:hover:not(:disabled),
.submit-btn:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 15px 35px rgba(79, 70, 229, 0.4);
}

.next-btn:disabled,
.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
  box-shadow: none !important;
}

.action-btn {
  min-width: 200px;
}

.action-btn.primary {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.btn-icon {
  font-size: 1.2em;
}

/* Индикатор вопросов */
.questions-indicator {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: center;
  padding: 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 15px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.question-dot {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f3f4f6;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-weight: 600;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.3s ease;
}

.question-dot:hover {
  background: #e5e7eb;
  transform: translateY(-2px);
}

.question-dot.current {
  background: #4f46e5;
  border-color: #4f46e5;
  color: white;
  transform: scale(1.1);
}

.question-dot.answered {
  background: #e0e7ff;
  border-color: #4f46e5;
  color: #4f46e5;
}

.question-dot.pending {
  opacity: 0.7;
}

/* Адаптивность */
@media (max-width: 768px) {
  .content-wrapper {
    padding: 15px;
  }
  
  .quiz-header {
    flex-direction: column;
    align-items: stretch;
    gap: 15px;
  }
  
  .exit-btn {
    align-self: flex-start;
  }
  
  .progress {
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
  }
  
  .question-card {
    padding: 20px;
  }
  
  .question-text {
    font-size: 1.2rem;
  }
  
  .answer-text {
    font-size: 1rem;
  }
  
  .navigation-buttons {
    flex-direction: column;
    gap: 15px;
  }
  
  .nav-btn,
  .next-btn,
  .submit-btn {
    width: 100%;
  }
  
  .score-display {
    flex-direction: column;
    gap: 30px;
  }
  
  .results-actions {
    flex-direction: column;
  }
}

@media (max-width: 480px) {
  .question-card {
    padding: 15px;
  }
  
  .answer-option {
    padding: 15px;
  }
  
  .question-dot {
    width: 35px;
    height: 35px;
    font-size: 0.9rem;
  }
  
  .score-circle {
    width: 120px;
    height: 120px;
  }
  
  .score-value {
    font-size: 2rem;
  }
}
</style>