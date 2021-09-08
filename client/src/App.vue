<script setup lang="ts">
import { createApp } from 'vue'
import PrimeVue from 'primevue/config';
import Main from './components/Main.vue'
import 'bulma/bulma.sass'

const App = createApp({})
App.use(PrimeVue);
</script>

<template>
  <div class="container base mx-auto">
    <div class="mx-3">
      <Main host="http://localhost:9081" perPage="10" />
    </div>
  </div>
</template>

<style lang="scss">
@import "./colors.scss";

body {
  color: $text-primary-color;
}

html {
  background-color: $background-color;
  height: 100%;
}
</style>
