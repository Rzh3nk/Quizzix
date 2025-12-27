
import { createApp } from 'vue'
import App from './App.vue'
import router from './router/router' // Маршрутизация


const app = createApp(App); // Создание экземпляра приложения
app.use(router); // Подключение маршрутизации
app.mount('#app')

