<script setup lang="ts">
import Main from './components/Main.vue'
import 'bulma/bulma.sass'
</script>

<template>
  <div class="container base mx-auto">
    <div class="mx-3">
      <Main perPage="10" />
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
