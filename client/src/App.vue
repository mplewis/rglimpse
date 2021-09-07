<script setup lang="ts">
import { createApp } from 'vue'
import PrimeVue from 'primevue/config';
import Main from './components/Main.vue'
import 'bulma/bulma.sass'

const App = createApp({})
App.use(PrimeVue);
</script>

<template>
  <div class="container mx-auto my-5">
    <Main host="http://localhost:9081" perPage="10" />
  </div>
</template>
