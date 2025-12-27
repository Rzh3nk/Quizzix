<template>
  <div id="app">
    
    <header>
      <h1>Quizzix</h1>
    </header>
    <main>
      <div v-if="loading">Загрузка квизов...</div>
      <ul v-else class="quiz-list">
        <li v-for="quiz in quizzes" :key="quiz.ID" class="quiz-item" @click="selectQuiz(quiz)">
          {{ quiz.title }}
        </li>
      </ul>
      <section v-if="selectedQuiz" class="questions">
        <h2>Вопросы квиза: {{ selectedQuiz.title }}</h2>
        <div v-if="questionsLoading" class="loading">Загрузка вопросов...</div>
        <ul v-else class="question">
          <li v-for="question in questions"
            :key="question.ID"
            class="question-item">
            {{ question.text }}
          </li>
        </ul>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'


const quizzes = ref([])
const loading = ref(true)

const selectedQuiz = ref(null)
const questions = ref([])
const questionsLoading = ref(false)
//подгрузка квизов
const fetchQuizzes = async () => {
  try {
    const response = await fetch('/api/quizzes')
    quizzes.value = await response.json()
  } catch (error) {
    console.error('Ошибка загрузки квизов:', error)
  } finally {
    loading.value = false
  }
}
//если выбрали квиз, надо лучше сделать чтобы переходило на другую страницу и там уже вопросы были
const selectQuiz = async (quiz) => {
  selectedQuiz.value = quiz
  questionsLoading.value = true
  questions.value = []

  try {
    const response = await fetch(`/api/quizzes/${quiz.ID}/questions`)
    questions.value = await response.json()
  } catch (error) {
    console.error('Ошибка загрузки вопросов:', error)
  } finally {
    questionsLoading.value = false
  }
}
//при загрузке вызывает
onMounted(fetchQuizzes)
</script>



<!-- Немного стили чтоб были -->
<style>
#app {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}
header h1 {
  color: #42b883;
  text-align: center;
  margin-bottom: 30px;
}
.quiz-list {
  list-style: none;
  padding: 0;
}
.quiz-item {
  padding: 15px 20px;
  background: blue;
  border-radius: 8px;
  margin-bottom: 10px;
  border-left: 4px solid #42b883;
  font-size: 18px;
  cursor: pointer;
}
.questions {
  margin-top: 30px;
}
.question-item {
  padding: 10px 15px;
  background: #457bdf;
  border-radius: 6px;
  margin-bottom: 8px;
  border: 1px solid #ddd;
}
</style>