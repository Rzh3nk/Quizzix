<template>
  <div id="app">
    <header>
      <h1>Quizzix</h1>
    </header>
    <main>
      <div v-if="loading" class="loading">Загрузка квизов...</div>
      <ul v-else class="quiz-list">
        <li v-for="quiz in quizzes" :key="quiz.id" class="quiz-item">
          {{ quiz.name }}
        </li>
      </ul>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const quizzes = ref([])
const loading = ref(true)

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

onMounted(fetchQuizzes)
</script>

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
}
.loading {
  text-align: center;
  padding: 40px;
  font-size: 18px;
  color: #666;
}
</style>
