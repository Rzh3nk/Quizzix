<template>
  <div id="app">
    <Register />
    
  </div>
</template>

<script setup>

import Register from './views/Register.vue';

//import HomeView from '@/views/HomeView.vue'

//import AppHeader from './components/AppHeader.vue'

</script>

<script>
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'App',
  setup() {
    // Здесь можно добавить реактивные данные или методы, если они понадобятся
    return {};
  }
});
</script>
