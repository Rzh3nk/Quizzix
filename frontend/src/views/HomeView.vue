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

/* Гарантируем полный экран БЕЗ лишнего flex */
.home-page {
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
  min-height: 100%; /* Вместо height: 100% */
}

/* Интересный фон с градиентом и узором */
.background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: 
    /* Тёмно-фиолетовый градиент */
    linear-gradient(135deg, #1e1b4b 0%, #312e81 25%, #4f46e5 50%, #7c3aed 75%, #a78bfa 100%),
    /* Сетчатый узор */
    repeating-linear-gradient(
      45deg,
      transparent,
      transparent 20px,
      rgba(255, 255, 255, 0.03) 20px,
      rgba(255, 255, 255, 0.03) 40px
    ),
    /* Точечный узор */
    radial-gradient(
      circle at 20% 80%,
      rgba(120, 119, 198, 0.1) 0%,
      transparent 50%
    ),
    radial-gradient(
      circle at 80% 20%,
      rgba(199, 210, 254, 0.1) 0%,
      transparent 50%
    );
  z-index: -1;
}

/* Центральный контейнер - УПРОЩАЕМ центрирование */
.content-wrapper {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
  max-width: 600px;
  margin: 0 auto;
  text-align: center;
  /* Убираем justify-content: center, чтобы не обрезало сверху */
}

/* Логотип и название на одной строке */
.logo-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin: 0 0 10px 0; /* Добавляем отступ сверху */
  width: 100%;
}

.logo-circle {
  width: 100px;
  height: 100px;
  background: white;
  border-radius: 25px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 3.5rem;
  font-weight: 900;
  color: #4f46e5;
  box-shadow: 
    0 10px 30px rgba(0, 0, 0, 0.3),
    0 0 0 2px rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
  flex-shrink: 0; /* Запрещаем сжатие логотипа */
}

.logo-circle:hover {
  transform: scale(1.05);
  box-shadow: 
    0 15px 40px rgba(0, 0, 0, 0.4),
    0 0 0 3px rgba(255, 255, 255, 0.3);
}

h1 {
  font-size: 3.2rem;
  color: white;
  margin: 0;
  padding: 0;
  font-weight: 800;
  text-shadow: 
    0 2px 10px rgba(0, 0, 0, 0.3),
    0 4px 20px rgba(0, 0, 0, 0.2);
  background-clip: text;
  line-height: 1.1; /* Устанавливаем line-height */
  display: block; /* Гарантируем блочное отображение */
}

/* Подзаголовок */
.subtitle {
  font-size: 1.3rem;
  color: rgba(255, 255, 255, 0.95);
  margin: 0 0 20px 0;
  font-weight: 500;
  text-align: center;
  width: 100%;
}

/* Описание */
.description {
  text-align: center;
  margin-bottom: 30px;
  width: 100%;
}

.description h2 {
  font-size: 1.8rem;
  color: white;
  margin: 0 0 15px 0;
  line-height: 1.3;
  font-weight: 600;
  text-align: center;
  display: block;
}

.description p {
  font-size: 1.1rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0 auto;
  line-height: 1.6;
  max-width: 500px;
  text-align: center;
  display: block;
}

/* Кнопки */
.buttons {
  width: 100%;
  max-width: 450px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  margin-bottom: 40px;
  margin-top: auto; /* Пушит кнопки вниз при необходимости */
}

.btn {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 22px 25px;
  border: none;
  border-radius: 18px;
  text-align: left;
  cursor: pointer;
  transition: all 0.3s ease;
  width: 100%;
  max-width: 400px;
}

.btn:hover {
  transform: translateY(-3px);
}

.icon {
  font-size: 1.8rem;
  flex-shrink: 0;
}

.text {
  flex: 1;
  text-align: left;
}

.text strong {
  display: block;
  font-size: 1.2rem;
  margin-bottom: 6px;
  font-weight: 700;
}

.text small {
  display: block;
  font-size: 0.9rem;
  opacity: 0.9;
  font-weight: 400;
}

.btn-login {
  background: white;
  color: #4f46e5;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.25);
}

.btn-login:hover {
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.35);
}

.btn-guest {
  background: rgba(255, 255, 255, 0.15);
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.25);
  backdrop-filter: blur(10px);
}

.btn-guest:hover {
  background: rgba(255, 255, 255, 0.25);
  border-color: rgba(255, 255, 255, 0.4);
}

.btn-register {
  background: rgba(255, 255, 255, 0.05);
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
}

.btn-register:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.5);
}

/* Адаптивность */
@media (max-width: 768px) {
  .content-wrapper {
    padding: 30px 15px;
    max-width: 500px;
  }
  
  .logo-row {
    gap: 15px;
    margin-top: 20px; /* Меньший отступ на мобильных */
  }
  
  .logo-circle {
    width: 80px;
    height: 80px;
    font-size: 2.8rem;
  }
  
  h1 {
    font-size: 2.5rem;
  }
  
  .subtitle {
    font-size: 1.1rem;
    margin-bottom: 30px;
  }
  
  .description h2 {
    font-size: 1.5rem;
  }
  
  .description p {
    font-size: 1rem;
    max-width: 400px;
  }
  
  .buttons {
    max-width: 400px;
  }
  
  .btn {
    padding: 18px 20px;
    max-width: 350px;
  }
}

@media (max-width: 480px) {
  .logo-row {
    flex-direction: column;
    gap: 10px;
    text-align: center;
    margin-top: 30px;
  }
  
  .logo-circle {
    width: 70px;
    height: 70px;
    font-size: 2.2rem;
  }
  
  h1 {
    font-size: 2.2rem;
    text-align: center;
  }
  
  .subtitle {
    font-size: 1rem;
    text-align: center;
  }
  
  .description h2 {
    font-size: 1.3rem;
    text-align: center;
  }
  
  .btn {
    padding: 16px 18px;
    max-width: 100%;
  }
  
  .icon {
    font-size: 1.5rem;
  }
  
  .text strong {
    font-size: 1.1rem;
  }
}

@media (max-width: 360px) {
  .content-wrapper {
    padding: 20px 10px;
  }
  
  h1 {
    font-size: 2rem;
  }
  
  .description h2 {
    font-size: 1.2rem;
  }
  
  .buttons {
    gap: 10px;
  }
}

/* ВАЖНО: Убираем лишнее центрирование из html/body */
html, body {
  height: 100%;
  margin: 0;
  padding: 0;
  /* Убираем display: flex чтобы не обрезало заголовок */
}

#app {
  height: 100%;
  /* Убираем display: flex чтобы не обрезало заголовок */
}

/* Альтернативный способ центрирования если нужно */
.alternative-centering {
  /* Если всё ещё есть проблемы с центрированием */
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
</style>