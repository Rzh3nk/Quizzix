<template>
  <div class="quiz-play-page">
    <div class="background"></div>

    <div class="content-wrapper">
      <!--Прогресс-->
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

      <div class="main-container">
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>Загрузка вопросов...</p>
        </div>

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
            <div class="question-header">
              <div class="question-meta">
                <span class="question-number">Вопрос {{ currentQuestionNumber }}</span>
              </div>
            </div>

            <!-- Текст вопроса -->
            <div class="question-content">
              <h2 class="question-text">{{ currentQuestion.text || currentQuestion.Text }}</h2>
              
              <div v-if="currentQuestion.img || currentQuestion.ImgURL" class="question-image">
                <img :src="getImageUrl(currentQuestion.img || currentQuestion.ImgURL)" :alt="currentQuestion.text">
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
                    <div v-if="isMultipleChoice" class="checkbox">
                      <div class="checkbox-inner" :class="{ 'checked': isAnswerSelected(answer) }"></div>
                    </div>
                    <div v-else class="radio">
                      <div class="radio-inner" :class="{ 'checked': isAnswerSelected(answer) }"></div>
                    </div>
                  </div>
                  
                  <div class="answer-content">
                    <div class="answer-text">{{ answer.text || answer.Text }}</div>
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

const quiz = ref(null)
const questions = ref([])
const loading = ref(true)
const error = ref(null)
const currentQuestionIndex = ref(0)
const userAnswers = ref({}) 
const isSubmitting = ref(false)
const showResults = ref(false)
const resultData = ref({
  score: 0,
  total: 0,
  percent: 0,
  details: []
})

const quizId = computed(() => route.params.id)

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

const correctAnswersCount = computed(() => {
  if (!currentQuestion.value?.answers) return 0
  return currentQuestion.value.answers.filter(answer => 
    answer.is_correct === true
  ).length
})

const isMultipleChoice = computed(() => {
  return correctAnswersCount.value > 1
})

const hasSelectedAnswers = computed(() => {
  if (!currentQuestion.value) return false
  
const questionId = currentQuestion.value.id || currentQuestion.value.ID
return Array.isArray(userAnswers.value[questionId]) && 
         userAnswers.value[questionId].length > 0
})


const fetchQuiz = async () => {
  try {
    loading.value = true
    error.value = null
    
    if (!quizId.value) {
      throw new Error('ID квиза не указан')
    }
    
    const token = localStorage.getItem('token')
    const headers = {
      'Content-Type': 'application/json'
    }
    
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    let quizUrl = `/api/quizzes/${quizId.value}`
    const quizResponse = await fetch(quizUrl, { headers })

    if (!quizResponse.ok) {
      if (quizResponse.status === 404) {
        throw new Error('Квиз не найден')
      }
      throw new Error(`Ошибка ${quizResponse.status}`)
    }

    const quizData = await quizResponse.json()

    quiz.value = {
      ID: quizData.id || quizData.ID,
      Title: quizData.title || quizData.Title,
      Description: quizData.description || quizData.Description
    }

    let questionsUrl = `/api/quizzes/${quizId.value}/questions`
    const questionsResponse = await fetch(questionsUrl, { headers })

    if (!questionsResponse.ok) {
      throw new Error('Вопросы не найдены')
    }

    const questionsData = await questionsResponse.json()
    questions.value = questionsData || []

    if (questions.value.length === 0) {
      throw new Error('В этом квизе нет вопросов')
    }

    initializeAnswers()
  } catch (err) {
    console.error('Ошибка загрузки квиза:', err)
    error.value = err.message || 'Не удалось загрузить вопросы'
  } finally {
    loading.value = false
  }
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
  
  if (isMultipleChoice.value) {
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
    
    // Получаем ID пользователя
    const userId = getUserId()
    
    // Подготавливаем данные для отправки
    const payload = {
      user_id: userId,
      answers: userAnswers.value
    }
    
    console.log('Отправляемые ответы:', payload)
    
    // Отправляем ответы на сервер
    let apiUrl = `/api/quizzes/${quizId.value}/submit`
   
    
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
  const userIdStr = localStorage.getItem('user_id')
  return userIdStr ? parseInt(userIdStr) : 0
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

onMounted(() => {
  if (quizId.value) {
    fetchQuiz()
  }
})
</script>

<style scoped>
@import '@/assets/que.css';
</style>