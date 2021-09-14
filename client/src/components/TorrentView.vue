<script setup lang="ts">
import { Torrent } from '../type'
import { useStore } from '../store'
import { prettyDate, prettySize } from '../logic/pretty'

const store = useStore()

defineProps<{ torrent: Torrent }>()

function eta(t: Torrent): Date {
  const remainingBytes = t.size - t.completed_bytes
  const remainingSecs = remainingBytes / t.down_rate  // down rate is in bytes per sec
  return new Date(Date.now() + remainingSecs * 1000)
}
</script>

<template>
  <div
    :class="{
      'torrent-card': true,
      'mb-3': true,
      'p-3': true,
      'complete': torrent.completed,
      'accessible-colors': store.state.accessibleColors,
    }"
  >
    <div v-if="torrent.completed">
      <div class="name has-text-weight-bold">{{ torrent.name }}</div>
      <div class="my-2 status-text">
        <progress :value="1" :max="1" />
      </div>
      <div class="columns is-mobile status-text">
        <div class="column">▲ {{ prettySize(torrent.up_rate) }}/s</div>
        <div class="column">{{ prettySize(torrent.completed_bytes) }}</div>
        <div class="column ratio">◕ {{ torrent.ratio }}</div>
      </div>
    </div>

    <div v-else>
      <div class="name has-text-weight-bold mb-2">{{ torrent.name }}</div>
      <div class="flex status-text">
        <div>{{ prettySize(torrent.size - torrent.completed_bytes) }} left</div>
        <span :title="eta(torrent).toString()">
          ETA:
          <span class="primary-text" v-if="torrent.down_rate > 0">{{ prettyDate(eta(torrent)) }}</span>
          <span class="primary-text" v-else>∞</span>
        </span>
      </div>
      <div class="my-2">
        <progress :value="torrent.completed_bytes" :max="torrent.size" />
      </div>
      <div class="columns is-mobile status-text">
        <div class="column">▼ {{ prettySize(torrent.down_rate) }}/s</div>
        <div class="column">{{ prettySize(torrent.size) }}</div>
        <div class="column ratio">◕ {{ torrent.ratio }}</div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
@use "sass:math";
@import "../colors.scss";

.column:nth-child(2) {
  text-align: center;
}
.column:nth-child(3) {
  text-align: right;
}

.torrent-card {
  border-radius: 6px;
  background-color: $torrent-unfinished-color;
  &.accessible-colors {
    background-color: $torrent-unfinished-color-accessible;
  }

  &.complete {
    background-color: $torrent-finished-color;
    &.accessible-colors {
      background-color: $torrent-finished-color-accessible;
    }
  }
}

.name {
  word-break: break-word;
}

$progress-height: 8px;

.progress-container {
  justify-content: left;
  align-items: center;
  display: flex;
}
progress,
progress[role] {
  height: $progress-height;
  width: 100%;
  background: $progress-background-color;
  border-radius: math.div($progress-height, 2);
  border: none;
  display: block;
}
progress[value]::-webkit-progress-bar {
  background: $progress-background-color;
  border-radius: math.div($progress-height, 2);
}
progress[value]::-moz-progress-bar {
  background: $progress-foreground-color;
  border-radius: math.div($progress-height, 2);
}
progress[value]::-webkit-progress-value {
  background: $progress-foreground-color;
  border-radius: math.div($progress-height, 2);
}
</style>
