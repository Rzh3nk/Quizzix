<template>
  <div class="create-quiz-page">
    <!-- Фон -->
    <div class="background"></div>
<Header/>
    <!-- Контент -->
    <div class="content-wrapper">
      <!-- Хлебные крошки -->
      <div class="breadcrumbs">
        <router-link to="/" class="breadcrumb-link">
          <span class="breadcrumb-icon">🏠</span>
          <span>Главная</span>
        </router-link>
        <span class="breadcrumb-separator">›</span>
        <span class="breadcrumb-current">Создать квиз</span>
      </div>

      <!-- Основной контейнер -->
      <div class="quiz-form-container">
        <div class="form-header">
          <h1>Создать новый квиз</h1>
          <p class="form-subtitle">Заполните информацию о квизе и добавьте вопросы</p>
        </div>

        <!-- Основная информация о квизе -->
        <div class="quiz-info-section">
          <h2 class="section-title">
            <span class="section-icon">📝</span>
            Основная информация
          </h2>
          
          <div class="fields-grid">
            <div class="field">
              <label class="field-label">
                <span class="label-icon">🏷️</span>
                Название квиза
                <span class="required">*</span>
              </label>
              <input 
                v-model="title" 
                type="text" 
                placeholder="Введите название квиза"
                class="input-field"
                :class="{ 'error': titleError }"
              />
              <div v-if="titleError" class="error-message">{{ titleError }}</div>
            </div>

            <div class="field">
              <label class="field-label">
                <span class="label-icon">📚</span>
                Категория
                <span class="required">*</span>
              </label>
              <div class="select-wrapper">
                <select 
                  v-model.number="categoryId" 
                  class="select-field"
                  :class="{ 'error': categoryError }"
                >
                  <option disabled value="0">Выберите категорию</option>
                  <option 
                    v-for="cat in categories" 
                    :key="cat.ID" 
                    :value="cat.ID"
                    class="option"
                  >
                    {{ cat.name }}
                  </option>
                </select>
              </div>
              <div v-if="categoryError" class="error-message">{{ categoryError }}</div>
            </div>
            <div class="field">
              <label class="field-label">
                <span class="label-icon">🖼️</span>
                Картинка квиза
                <span class="required">*</span>
              </label>
              <input 
                v-model="imgUrl" 
                type="url" 
                placeholder="https://example.com/image.jpg"
                class="input-field"
                :class="{ 'error': imgUrlError }"
              />
              <div class="hint">Ссылка на изображение</div>
              <div v-if="imgUrlError" class="error-message">{{ imgUrlError }}</div>
            </div>
            <div class="field full-width">
              <label class="field-label">
                <span class="label-icon">📄</span>
                Описание квиза
              </label>
              <textarea 
                v-model="description" 
                rows="3" 
                placeholder="Опишите содержание квиза, его тематику и особенности"
                class="textarea-field"
              />
              <div class="char-counter">{{ description.length }}/500</div>
            </div>
          </div>
        </div>

        <!-- Список вопросов -->
        <div class="questions-section">
          <div class="section-header">
            <h2 class="section-title">
              <span class="section-icon">❓</span>
              Вопросы
              <span class="questions-count">({{ questions.length }})</span>
            </h2>
          </div>

          <!-- Сообщение если нет вопросов -->
          <div v-if="questions.length === 0" class="empty-questions">
            <div class="empty-icon">📝</div>
            <h3>Пока нет вопросов</h3>
            <p>Добавьте первый вопрос для вашего квиза</p>
          </div>

          <!-- Список вопросов -->
          <div class="questions-list">
            <div
              v-for="(question, qIndex) in questions"
              :key="qIndex"
              class="question-card"
            >
              <div class="question-header">
                <div class="question-number">
                  <span class="number">Вопрос {{ qIndex + 1 }}</span>
                  <span v-if="question.text" class="question-preview">
                    {{ question.text.substring(0, 50) }}{{ question.text.length > 50 ? '...' : '' }}
                  </span>
                </div>
                <button 
                  type="button" 
                  @click="removeQuestion(qIndex)"
                  class="remove-question-btn"
                  :disabled="questions.length <= 1"
                >
                  <span class="remove-icon">🗑️</span>
                  Удалить вопрос
                </button>
              </div>

              <div class="question-content">
                <!-- Текст вопроса -->
                <div class="field">
                  <label class="field-label">
                    <span class="label-icon">💬</span>
                    Текст вопроса
                    <span class="required">*</span>
                  </label>
                  <input
                    v-model="question.text"
                    type="text"
                    placeholder="Введите текст вопроса"
                    class="input-field"
                    :class="{ 'error': questionError(qIndex) }"
                  />
                  <div v-if="questionError(qIndex)" class="error-message">
                    {{ questionError(qIndex) }}
                  </div>
                </div>
                <div class="field">
                  <label class="field-label">
                    <span class="label-icon">🖼️</span>
                    Картинка к вопросу
                  </label>
                  <input
                    v-model="question.img"
                    type="url"
                    placeholder="https://example.com/image.jpg"
                    class="input-field"
                  />
                  <div class="hint">Ссылка на изображение (опционально)</div>
                  
                  <!-- ✅ ПРЕДВАРИТЕЛЬНЫЙ ПРОСМОТР -->
                  <div v-if="question.img" class="image-preview">
                    <img 
                      :src="question.img" 
                      :alt="question.text.substring(0, 30) + '...'" 
                      @error="question.img = ''"
                    />
                  </div>
                </div>
                <!-- Ответы для вопроса -->
                <div class="answers-section">
                  <div class="answers-header">
                    <h3 class="answers-title">
                      <span class="answers-icon">📝</span>
                      Ответы
                      <span class="answers-count">({{ question.answers.length }})</span>
                    </h3>
                  </div>

                  <!-- Список ответов -->
                  <div class="answers-list">
                    <div
                      v-for="(answer, aIndex) in question.answers"
                      :key="aIndex"
                      class="answer-card"
                      :class="{ 'correct': answer.is_correct }"
                    >
                      <div class="answer-content">
                        <div class="answer-input">
                          <input
                            v-model="answer.text"
                            type="text"
                            placeholder="Текст ответа"
                            class="input-field"
                            :class="{ 'error': answerError(qIndex, aIndex) }"
                          />
                          <div v-if="answerError(qIndex, aIndex)" class="error-message">
                            {{ answerError(qIndex, aIndex) }}
                          </div>
                        </div>
                        
                        <div class="answer-controls">
                          <label class="correct-checkbox">
                            <input
                              type="checkbox"
                              v-model="answer.is_correct"
                              class="checkbox-input"
                            />
                            <span class="checkbox-custom"></span>
                            <span class="checkbox-label">
                              <span class="correct-icon">✓</span>
                              Правильный ответ
                            </span>
                          </label>
                          
                          <button 
                            type="button" 
                            @click="removeAnswer(qIndex, aIndex)"
                            class="remove-answer-btn"
                            :disabled="question.answers.length <= 2"
                          >
                            Удалить ответ
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Кнопка добавления ответа внизу секции ответов -->
                  <div class="add-answer-footer">
                    <button 
                      type="button" 
                      @click="addAnswer(qIndex)"
                      class="add-answer-btn"
                    >
                      <span class="btn-icon">+</span>
                      Добавить ответ
                    </button>
                  </div>

                  <!-- Сообщение о правильных ответах -->
                  <div v-if="getCorrectAnswersCount(qIndex) === 0" class="warning-message">
                    ⚠️ Необходимо отметить хотя бы один правильный ответ
                  </div>
                  <div v-else class="success-message">
                    ✓ Отмечено правильных ответов: {{ getCorrectAnswersCount(qIndex) }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Кнопка добавления вопроса внизу секции вопросов -->
          <div class="add-question-footer">
            <button 
              type="button" 
              @click="addQuestion" 
              class="add-question-btn"
            >
              <span class="btn-icon">+</span>
              Добавить новый вопрос
            </button>
          </div>
        </div>

        <!-- Действия -->
        <div class="actions-section">
          <div class="actions-info">
            <div class="stats">
              <div class="stat">
                <span class="stat-icon">📝</span>
                <span class="stat-label">Вопросов:</span>
                <span class="stat-value">{{ questions.length }}</span>
                
              </div>
              <div class="stat">
                <span class="stat-icon">❓</span>
                <span class="stat-label">Ответов:</span>
                <span class="stat-value">{{ totalAnswers }}</span>
                
              </div>
            </div>
          </div>

          <div class="action-buttons">
            
            <button 
              type="button" 
              @click="createQuiz" 
              class="primary-btn"
              :disabled="isSubmitting || !isFormValid"
            >
              <span v-if="isSubmitting" class="spinner"></span>
              <span v-else class="btn-icon">🚀</span>
              {{ isSubmitting ? 'Создание...' : 'Создать квиз' }}
            </button>
          </div>
        </div>

        <!-- Сообщения -->
        <div v-if="error" class="error-container">
          <span class="error-icon">⚠️</span>
          <p class="error-text">{{ error }}</p>
          <button @click="error = ''" class="close-error-btn">
            ✕
          </button>
        </div>

        <div v-if="success" class="success-container">
          <span class="success-icon">🎉</span>
          <div class="success-content">
            <h3 class="success-title">Квиз успешно создан!</h3>
            <p class="success-text">{{ success }}</p>
            <div class="success-actions">
              <router-link :to="`/quiz/${createdQuizId}`" class="success-btn">
                <span class="btn-icon">👁️</span>
                Посмотреть квиз
              </router-link>
              <button @click="resetForm" class="success-btn secondary">
                <span class="btn-icon">➕</span>
                Создать еще один
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Кнопка назад -->
      <button @click="goBack" class="back-btn">
        <span class="back-icon">←</span>
        <span>Назад</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
import Header from '@/components/Header.vue'
// Основные данные квиза
const title = ref('')
const description = ref('')
const categoryId = ref(0)
const difficulty = ref('medium')
const timeLimit = ref(30)
const imgUrl = ref('')

const imgUrlError = computed(() => {
  if (!imgUrl.value.trim()) return 'Картинка обязательна'
  if (!imgUrl.value.match(/^https?:\/\/.+\.(jpg|jpeg|png|webp)$/i)) {
    return 'Введите корректную ссылку на изображение'
  }
  return ''
})
// Списки
const categories = ref([])
const questions = ref([
  {
    text: '',
    img: '',
    answers: [
      { text: '', is_correct: false },
      { text: '', is_correct: true },
    ],
  },
])

// Состояния
const error = ref('')
const success = ref('')
const isSubmitting = ref(false)
const createdQuizId = ref(null)

// Опции сложности
const difficultyOptions = [
  { value: 'easy', label: 'Лёгкий', icon: '😊' },
  { value: 'medium', label: 'Средний', icon: '😐' },
  { value: 'hard', label: 'Сложный', icon: '😓' }
]

// Компьютеды
const totalAnswers = computed(() => {
  return questions.value.reduce((sum, q) => sum + q.answers.length, 0)
})

const titleError = computed(() => {
  if (!title.value.trim()) return 'Введите название квиза'
  if (title.value.length < 3) return 'Название слишком короткое'
  if (title.value.length > 100) return 'Название слишком длинное'
  return ''
})

const categoryError = computed(() => {
  if (!categoryId.value) return 'Выберите категорию'
  return ''
})

const isFormValid = computed(() => {
  if (titleError.value) return false
  if (categoryError.value) return false
  if (imgUrlError.value) return false
  if (questions.value.length === 0) return false
  
  // Проверка всех вопросов
  for (const question of questions.value) {
    if (!question.text.trim()) return false
    if (question.answers.length < 2) return false
    if (!question.answers.some(a => a.is_correct)) return false
    if (question.answers.some(a => !a.text.trim())) return false
  }
  
  return true
})

// Методы валидации
const questionError = (qIndex) => {
  const question = questions.value[qIndex]
  if (!question.text.trim()) return 'Введите текст вопроса'
  if (question.text.length < 5) return 'Вопрос слишком короткий'
  return ''
}

const answerError = (qIndex, aIndex) => {
  const answer = questions.value[qIndex].answers[aIndex]
  if (!answer.text.trim()) return 'Введите текст ответа'
  return ''
}

// Методы
const selectDifficulty = (diff) => {
  difficulty.value = diff
}

const getDifficultyText = (diff) => {
  const map = { easy: 'Лёгкий', medium: 'Средний', hard: 'Сложный' }
  return map[diff] || diff
}

const addQuestion = () => {
  questions.value.push({
    text: '',
    img:'',
    answers: [
      { text: '', is_correct: false },
      { text: '', is_correct: true },
    ],
  })
}

const removeQuestion = (index) => {
  if (questions.value.length > 1) {
    questions.value.splice(index, 1)
  }
}

const addAnswer = (qIndex) => {
  questions.value[qIndex].answers.push({
    text: '',
    is_correct: false,
  })
}

const removeAnswer = (qIndex, aIndex) => {
  if (questions.value[qIndex].answers.length > 2) {
    questions.value[qIndex].answers.splice(aIndex, 1)
  }
}

const getCorrectAnswersCount = (qIndex) => {
  return questions.value[qIndex].answers.filter(a => a.is_correct).length
}

// Загрузка категорий
const fetchCategories = async () => {
  try {
    const response = await fetch(`/api/categories`)
    if (response.ok) {
      categories.value = await response.json()
    }
  } catch (err) {
    console.error('Ошибка загрузки категорий:', err)
    // Мок данные для разработки
    if (import.meta.env.DEV) {
      categories.value = [
        { ID: 1, name: 'История' },
        { ID: 2, name: 'Наука' },
        { ID: 3, name: 'Литература' },
        { ID: 4, name: 'География' },
        { ID: 5, name: 'Музыка' },
        { ID: 6, name: 'Кино' },
        { ID: 7, name: 'Спорт' },
        { ID: 8, name: 'Технологии' }
      ]
    }
  }
}

// Создание квиза
const createQuiz = async () => {
  error.value = ''
  success.value = ''
const authorId = localStorage.getItem('user_id')  // "1"
  console.log('🔍 authorId:', authorId)
  if (!isFormValid.value) {
    error.value = 'Пожалуйста, заполните все обязательные поля правильно'
    return
  }
  
  isSubmitting.value = true
  
  try {
    const payload = {
      author_id: parseInt(authorId),
      title: title.value.trim(),
      description: description.value.trim(),
      img: imgUrl.value.trim(),
      category_id: categoryId.value,
      difficulty: difficulty.value,
      time_limit: timeLimit.value || 0,
      questions: questions.value.map((q) => ({
        text: q.text.trim(),
        img: q.img.trim(),
        answers: q.answers.map((a) => ({
          text: a.text.trim(),
          is_correct: a.is_correct,
        })),
      })),
    }
    
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/quizzes`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(payload)
    })
    
    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Ошибка ${response.status}: ${errorText}`)
    }
    
    const data = await response.json()
    createdQuizId.value = data.id
    
    success.value = `Квиз "${title.value}" успешно создан!`
    
  } catch (err) {
    console.error('Ошибка при создании квиза:', err)
    error.value = err.message || 'Не удалось создать квиз. Попробуйте еще раз.'
  } finally {
    isSubmitting.value = false
  }
}

// Дополнительные методы
const saveDraft = () => {
  // Сохранение в localStorage
  const draft = {
    title: title.value,
    description: description.value,
    categoryId: categoryId.value,
    difficulty: difficulty.value,
    timeLimit: timeLimit.value,
    questions: questions.value,
    savedAt: new Date().toISOString()
  }
  
  localStorage.setItem('quiz_draft', JSON.stringify(draft))
  alert('Черновик сохранён!')
}

const resetForm = () => {
  title.value = ''
  description.value = ''
  imgUrl.value = '' 
  categoryId.value = 0
  difficulty.value = 'medium'
  timeLimit.value = 30
  questions.value = [{
    text: '',
    answers: [
      { text: '', is_correct: false },
      { text: '', is_correct: true },
    ],
  }]
  error.value = ''
  success.value = ''
  createdQuizId.value = null
}

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/')
  }
}

// Хуки
onMounted(() => {
  fetchCategories()
  
  // Загрузка черновика из localStorage
  const savedDraft = localStorage.getItem('quiz_draft')
  if (savedDraft) {
    try {
      const draft = JSON.parse(savedDraft)
      title.value = draft.title
      description.value = draft.description
      categoryId.value = draft.categoryId
      difficulty.value = draft.difficulty || 'medium'
      timeLimit.value = draft.timeLimit || 30
      questions.value = draft.questions || [{
        text: '',
        answers: [
          { text: '', is_correct: false },
          { text: '', is_correct: true },
        ],
      }]
    } catch (e) {
      console.error('Ошибка загрузки черновика:', e)
    }
  }
})
</script>

<style scoped>
/* Основные стили страницы */
.create-quiz-page {
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
}

.breadcrumb-link {
  display: flex;
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

/* Контейнер формы */
.quiz-form-container {
  width: 100%;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.form-header {
  text-align: center;
  margin-bottom: 40px;
}

.form-header h1 {
  font-size: 2.5rem;
  margin-bottom: 10px;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.form-subtitle {
  color: #666;
  font-size: 1.1rem;
  opacity: 0.8;
}

/* Секции */
.quiz-info-section,
.questions-section,
.actions-section {
  margin-bottom: 40px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1.5rem;
  margin-bottom: 25px;
  color: #333;
  padding-bottom: 10px;
  border-bottom: 2px solid #f3f4f6;
}

.section-icon {
  font-size: 1.3em;
}

.questions-count,
.answers-count {
  color: #666;
  font-size: 0.9em;
  margin-left: 5px;
  font-weight: normal;
}

/* Поля формы */
.fields-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 100px;
  margin-bottom: 30px;
}

.field {
  text-align: left;
}

.field.full-width {
  grid-column: 1 / -1;
}

.field-label {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  font-weight: 600;
  color: #444;
  font-size: 0.95rem;
}

.label-icon {
  font-size: 1.1em;
}

.required {
  color: #ef4444;
  margin-left: 2px;
}

.input-field,
.textarea-field,
.select-field {
  width: 97%;
  padding: 10px;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  font-size: 1rem;
  background: white;
  color: #1f2937;
  transition: all 0.3s ease;
  font-family: inherit;
}

.input-field:focus,
.textarea-field:focus,
.select-field:focus {
  outline: none;
  border-color: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79, 70, 229, 0.1);
}

.input-field.error,
.select-field.error {
  border-color: #ef4444;
}

.textarea-field {
  resize: vertical;
  min-height: 100px;
  line-height: 1.5;
}

.select-wrapper {
  position: relative;
}

.select-arrow {
  position: absolute;
  right: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: #6b7280;
  pointer-events: none;
  font-size: 0.8em;
}

.char-counter {
  text-align: right;
  font-size: 0.85rem;
  color: #6b7280;
  margin-top: 5px;
}

.hint {
  font-size: 0.85rem;
  color: #6b7280;
  margin-top: 5px;
  font-style: italic;
}

/* Кнопки сложности */
.difficulty-buttons {
  display: flex;
  gap: 10px;
}

.diff-btn {
  flex: 1;
  padding: 12px 15px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  background: white;
  color: #4b5563;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 0.9rem;
}

.diff-btn:hover {
  border-color: #d1d5db;
  transform: translateY(-2px);
}

.diff-btn.selected {
  color: white;
  border-color: transparent;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.diff-btn.easy.selected {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
}

.diff-btn.medium.selected {
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
}

.diff-btn.hard.selected {
  background: linear-gradient(135deg, #ef4444 0%, #f87171 100%);
}

.diff-icon {
  font-size: 1.2em;
}

/* Раздел вопросов */
.questions-section {
  background: #f9fafb;
  border-radius: 15px;
  padding: 30px;
  margin-top: 30px;
}

.section-header {
  margin-bottom: 25px;
}

.empty-questions {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 15px;
  border: 2px dashed #d1d5db;
  margin-bottom: 30px;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 20px;
  color: #9ca3af;
  opacity: 0.7;
}

.empty-questions h3 {
  color: #4b5563;
  margin-bottom: 10px;
  font-size: 1.3rem;
}

.empty-questions p {
  color: #6b7280;
  font-size: 1rem;
}

/* Карточка вопроса */
.questions-list {
  display: flex;
  flex-direction: column;
  gap: 25px;
  margin-bottom: 30px;
}

.question-card {
  background: white;
  border-radius: 15px;
  border: 2px solid #e5e7eb;
  padding: 25px;
  transition: all 0.3s ease;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
}

.question-card:hover {
  border-color: #cbd5e1;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
}

.question-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 2px solid #f3f4f6;
}

.question-number {
  display: flex;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}

.number {
  font-weight: 700;
  color: #4f46e5;
  font-size: 1.2rem;
  background: #eef2ff;
  padding: 8px 15px;
  border-radius: 10px;
}

.question-preview {
  color: #6b7280;
  font-size: 0.95rem;
  max-width: 400px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.remove-question-btn {
  padding: 10px 20px;
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
  border-radius: 10px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
}

.remove-question-btn:hover:not(:disabled) {
  background: #fee2e2;
  border-color: #fca5a5;
  transform: translateY(-1px);
}

.remove-question-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.remove-icon {
  font-size: 1.1em;
}

.question-content {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

/* Раздел ответов */
.answers-section {
  background: #f9fafb;
  border-radius: 12px;
  padding: 25px;
  border: 1px solid #e5e7eb;
}

.answers-header {
  margin-bottom: 20px;
}

.answers-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 1.2rem;
  color: #444;
  margin: 0;
  padding-bottom: 10px;
  border-bottom: 1px solid #e5e7eb;
}

.answers-icon {
  font-size: 1.2em;
}

.answers-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-bottom: 20px;
}

.answer-card {
  background: white;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  padding: 20px;
  transition: all 0.3s ease;
}

.answer-card.correct {
  border-color: #10b981;
  background: linear-gradient(135deg, #f0fdf4 0%, #dcfce7 100%);
  box-shadow: 0 3px 10px rgba(16, 185, 129, 0.1);
}

.answer-content {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.answer-input {
  flex: 1;
}

.answer-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  flex-wrap: wrap;
}

.correct-checkbox {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  flex: 1;
}

.checkbox-input {
  display: none;
}

.checkbox-custom {
  width: 22px;
  height: 22px;
  border: 2px solid #d1d5db;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  flex-shrink: 0;
}

.checkbox-input:checked + .checkbox-custom {
  background: #10b981;
  border-color: #10b981;
}

.checkbox-input:checked + .checkbox-custom::after {
  content: '✓';
  color: white;
  font-weight: bold;
  font-size: 0.9em;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: #374151;
  font-size: 0.95rem;
}

.correct-icon {
  color: #10b981;
  font-weight: bold;
  font-size: 1.1em;
}

.remove-answer-btn {
  padding: 8px 16px;
  background: #f3f4f6;
  color: #6b7280;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.9rem;
}

.remove-answer-btn:hover:not(:disabled) {
  background: #ef4444;
  color: white;
  border-color: #ef4444;
  transform: translateY(-1px);
}

.remove-answer-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Кнопка добавления ответа внизу секции ответов */
.add-answer-footer {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #e5e7eb;
  text-align: center;
}

.add-answer-btn {
  padding: 12px 25px;
  background: white;
  color: #4f46e5;
  border: 2px solid #4f46e5;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 0.95rem;
}

.add-answer-btn:hover {
  background: #4f46e5;
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(79, 70, 229, 0.2);
}

/* Кнопка добавления вопроса внизу секции вопросов */
.add-question-footer {
  margin-top: 30px;
  padding-top: 30px;
  border-top: 2px solid #e5e7eb;
  text-align: center;
}

.add-question-btn {
  padding: 15px 35px;
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  box-shadow: 0 5px 15px rgba(79, 70, 229, 0.3);
}

.add-question-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 10px 25px rgba(79, 70, 229, 0.4);
}

.btn-icon {
  font-size: 1.2em;
}

/* Сообщения о правильных ответах */
.warning-message,
.success-message {
  margin-top: 20px;
  padding: 12px 15px;
  border-radius: 8px;
  font-weight: 500;
  text-align: left;
  font-size: 0.9rem;
}

.warning-message {
  background: #fef3c7;
  color: #92400e;
  border: 1px solid #fde68a;
}

.success-message {
  background: #d1fae5;
  color: #065f46;
  border: 1px solid #a7f3d0;
}

.error-message {
  color: #dc2626;
  font-size: 0.85rem;
  margin-top: 5px;
  font-weight: 500;
}

/* Раздел действий */
.actions-section {
  margin-top: 40px;
  padding-top: 30px;
  border-top: 2px solid #f3f4f6;
}

.actions-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  flex-wrap: wrap;
  gap: 20px;
}

.stats {
  display: flex;
  gap: 30px;
  flex-wrap: wrap;
}

.stat {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #f9fafb;
  padding: 10px 15px;
  border-radius: 10px;
  border: 1px solid #e5e7eb;
}

.stat-icon {
  font-size: 1.2em;
}

.stat-value {
  font-weight: 700;
  color: #4f46e5;
  font-size: 1.1rem;
}

.stat-label {
  color: #6b7280;
  font-size: 0.9rem;
}

.action-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  flex-wrap: wrap;
}

.primary-btn,
.secondary-btn {
  padding: 15px 35px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 10px;
  border: none;
  min-width: 200px;
  justify-content: center;
}

.primary-btn {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  color: white;
  box-shadow: 0 10px 25px rgba(79, 70, 229, 0.3);
}

.primary-btn:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 15px 35px rgba(79, 70, 229, 0.4);
}

.primary-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
}

.secondary-btn {
  background: white;
  color: #4f46e5;
  border: 2px solid #4f46e5;
}

.secondary-btn:hover:not(:disabled) {
  background: #f8fafc;
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(79, 70, 229, 0.1);
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-radius: 50%;
  border-top-color: white;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Сообщения об ошибках и успехах */
.error-container,
.success-container {
  margin-top: 30px;
  padding: 20px;
  border-radius: 12px;
  display: flex;
  align-items: flex-start;
  gap: 15px;
}

.error-container {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #dc2626;
}

.success-container {
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  color: #065f46;
}

.error-icon,
.success-icon {
  font-size: 1.5rem;
  flex-shrink: 0;
}

.error-text,
.success-content {
  flex: 1;
  text-align: left;
}

.success-title {
  font-size: 1.3rem;
  margin-bottom: 5px;
  font-weight: 600;
}

.success-text {
  margin-bottom: 15px;
  font-size: 0.95rem;
}

.success-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.success-btn {
  padding: 10px 20px;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
  border: none;
  cursor: pointer;
}

.success-btn {
  background: #10b981;
  color: white;
  border: 1px solid #10b981;
}

.success-btn.secondary {
  background: white;
  color: #10b981;
  border: 2px solid #10b981;
}

.success-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.close-error-btn {
  background: none;
  border: none;
  font-size: 1.2rem;
  color: #dc2626;
  cursor: pointer;
  padding: 0;
  opacity: 0.7;
  transition: opacity 0.2s;
  flex-shrink: 0;
}

.close-error-btn:hover {
  opacity: 1;
}

/* Кнопка назад */
.back-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 15px 30px;
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.2);
  color: white;
  border-radius: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 30px;
  backdrop-filter: blur(10px);
  font-size: 1rem;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.back-icon {
  font-size: 1.1em;
}

</style>