<script setup lang="ts">
import { ref, watchEffect } from 'vue'
import { useStore } from '../store'

const store = useStore()

const accessibleColors = ref<boolean>(store.state.accessibleColors)
watchEffect(() => store.commit('setAccessibleColors', accessibleColors.value))
</script>

<template>
  <teleport to="body">
    <div v-if="store.state.prefsOpen" class="settings p-4">
      <div class="button-holder">
        <button class="button is-primary is-small" @click="store.commit('togglePrefs')">&cross;</button>
      </div>
      <h1 class="has-text-weight-bold is-size-5 mb-2">Settings</h1>
      <div class="mt-1">
        <input type="checkbox" class="checkbox" id="accessible-color" v-model="accessibleColors" />
        <label for="accessible-color" class="ml-2">Use accessible (colorblind) colors</label>
      </div>
    </div>
  </teleport>
</template>

<style scoped lang="scss">
@import "../colors.scss";

.settings {
  position: fixed;
  background: $background-color;
  border-radius: 6px;
  box-shadow: 0 0 8px rgba(0, 0, 0, 1);
  z-index: 1;
  top: 20px;
  left: 20px;

  .button-holder {
    text-align: right;
    height: 0;
  }
}
</style>
