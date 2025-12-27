<template>
  <div class="home-page">
    <!-- Фон с узором -->
    <div class="background"></div>
    
    <!-- Контент -->
    <div class="content-wrapper">
      <!-- Логотип и название на одной строке -->
      <div class="logo-row">
        <div class="logo-circle">Q</div>
        <h1>Quizzix</h1>
      </div>
      
      <!-- Подзаголовок -->
      <p class="subtitle">Платформа для умных квизов</p>
      
      <!-- Описание -->
      <div class="description">
        <h2>Создавайте и проходите увлекательные квизы</h2>
        <p>Проверяйте знания, соревнуйтесь с друзьями, создавайте свои тесты</p>
      </div>
      
      <!-- Кнопки -->
      <div class="buttons">
        <button class="btn btn-login" @click="goToLogin">
          <span class="icon">🔑</span>
          <div class="text">
            <strong>Войти</strong>
            <small>Для сохранения прогресса</small>
          </div>
        </button>
        
        <button class="btn btn-guest" @click="continueAsGuest">
          <span class="icon">🚀</span>
          <div class="text">
            <strong>Продолжить без входа</strong>
            <small>Быстрый доступ</small>
          </div>
        </button>
        
        <button class="btn btn-register" @click="goToRegister">
          <span class="icon">✨</span>
          <div class="text">
            <strong>Создать аккаунт</strong>
            <small>Полный доступ</small>
          </div>
        </button>
      </div>
      
    </div>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'

export default {
  name: 'HomeView',
  
  setup() {
    const router = useRouter()
    
    const goToLogin = () => {
      console.log('Переход на страницу входа')
      alert('Страница входа будет доступна после создания компонента LoginView')
    }
    
    const goToRegister = () => {
      console.log('Переход на страницу регистрации')
      alert('Страница регистрации будет доступна после создания компонента RegisterView')
    }
    
    const continueAsGuest = () => {
      console.log('Продолжить как гость')
      alert('Гостевой режим будет доступен после создания компонента CategoriesView')
    }
    
    return {
      goToLogin,
      goToRegister,
      continueAsGuest
    }
  }
}
</script>

<style scoped>
/* Стили только для этого компонента */
@import '@/assets/home.css';
</style>
