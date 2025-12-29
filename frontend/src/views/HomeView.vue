<template>
  <div class="home-page">
    <div class="background"></div>
    <div class="content-wrapper">
      <div class="logo-row">
        <div class="logo-circle">Q</div>
        <h1>Quizzix</h1>
      </div>
      
      <p class="subtitle">Платформа для умных квизов</p>
      
      <div class="description">
        <h2>Создавайте и проходите увлекательные квизы</h2>
        <p>Проверяйте знания, соревнуйтесь с друзьями, создавайте свои тесты</p>
      </div>
      <!--Кнопки-->
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
      router.push('/login')
    }
    
    const goToRegister = () => {
      console.log('Переход на страницу регистрации')
      router.push('/register')
    }
    
    const continueAsGuest = () => {
      console.log('Продолжить как гость')
      router.push('/main')
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
@import '@/assets/home.css';
</style>
