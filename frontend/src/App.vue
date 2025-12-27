<template>
  <div id="app">
    <router-view />
  </div>
</template>


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
