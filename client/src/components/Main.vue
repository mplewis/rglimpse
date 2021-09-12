<script setup lang="ts">
import { ref, watchEffect } from 'vue'
import { prettyDate, prettySize } from '../logic/pretty'

const props = defineProps<{
  perPage: string,
}>()
const perPage = parseInt(props.perPage)

const loading = ref<boolean>(true)
const error = ref<any>(null)
const name = ref<string | null>(null)
const page = ref<number>(1)
const total = ref<number>(0)
const torrents = ref<Torrent[]>([])
const query = ref<string>('')
const queryRaw = ref<string>('')

const showSettings = ref<boolean>(false)
const accessibleColor = ref<boolean>(false)

type Torrent = {
  completed: boolean,
  completed_bytes: number,
  created: Date,
  down_rate: number,
  finished: Date,
  hash: string,
  label: string,
  name: string,
  path: string,
  ratio: number,
  size: number,
  started: Date,
  up_rate: number,
}

function eta(t: Torrent): Date {
  const remainingBytes = t.size - t.completed_bytes
  const remainingSecs = remainingBytes / t.down_rate  // down rate is in bytes per sec
  return new Date(Date.now() + remainingSecs * 1000)
}

function parse(raw: { [k: string]: any }): Torrent {
  return {
    completed: raw.completed,
    completed_bytes: raw.completed_bytes,
    created: new Date(raw.created),
    down_rate: raw.down_rate,
    finished: new Date(raw.finished),
    hash: raw.hash,
    label: raw.label,
    name: raw.name,
    path: raw.path,
    ratio: raw.ratio,
    size: raw.size,
    started: new Date(raw.started),
    up_rate: raw.up_rate,
  }
}

// TODO: auto refresh
async function fetchTorrents() {
  try {
    const count = parseInt(props.perPage)
    const offset = (page.value - 1) * count
    var path = `/torrents?offset=${offset}&count=${count}`
    if (query.value.length > 0) {
      path = `/torrents?offset=${offset}&count=${count}&query=${query.value}`
    }

    const resp = await fetch(path)
    const data = await resp.json()
    name.value = data.name
    total.value = data.total
    torrents.value = data.torrents.map(parse)
    error.value = null
  } catch (e) {
    console.log(e)
    error.value = e
  } finally {
    loading.value = false
  }
}

function diffPage(diff: number) {
  page.value += diff
}

watchEffect(() => {
  fetchTorrents()
})

setInterval(fetchTorrents, 5000)

// https://decipher.dev/30-seconds-of-typescript/docs/debounce/
const debounce = (fn: Function, ms = 300) => {
  let timeoutId: ReturnType<typeof setTimeout>;
  return function (this: any, ...args: any[]) {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => fn.apply(this, args), ms);
  };
};

const setQuery = debounce((q: string) => query.value = q, 300)
watchEffect(() => setQuery(queryRaw.value))
</script>

<template>
  <!--
  <div v-if="showSettings" class="settings p-4">
    <div class="button-holder">
      <button class="button is-primary is-small" @click="showSettings = false">&cross;</button>
    </div>
    <h1 class="has-text-weight-bold is-size-5 mb-2">Settings</h1>
    <div class="mt-1">
      <input type="checkbox" class="checkbox" id="accessible-color" v-model="accessibleColor" />
      <label for="accessible-color" class="ml-2">Use accessible (colorblind) colors</label>
    </div>
  </div>
  -->

  <div class="flex">
    <h1 class="app-name is-size-3 mb-3" @click="showSettings = true">rglimpse</h1>
    <div class="app-state has-text-right">
      <p v-if="loading">Loading...</p>
      <p v-else-if="name">{{ name }}</p>
    </div>
  </div>

  <div v-if="error">
    <p>Sorry, something went wrong:</p>
    <pre class="my-3"><code>{{ error.toString() }}</code></pre>
    <p>Try refreshing the page.</p>
  </div>

  <div v-else>
    <div class="mb-3">
      <input
        type="text"
        class="input search"
        placeholder="Search for torrents..."
        v-model="queryRaw"
      />
    </div>

    <div class="flex mb-4" v-if="total">
      <button class="button is-primary" @click="diffPage(-1)" :disabled="page <= 1">&laquo;</button>
      {{ (page - 1) * perPage + 1 }}–{{ Math.min(page * perPage, total) }} of {{ total }}
      <button
        class="button is-primary"
        @click="diffPage(1)"
        :disabled="page >= Math.ceil(total / perPage)"
      >&raquo;</button>
    </div>

    <div
      v-if="!loading"
      v-for="torrent in torrents"
      :key="torrent.hash"
      :class="{ 'torrent-card': true, 'mb-3': true, 'p-3': true, 'complete': torrent.completed }"
    >
      <div v-if="torrent.completed" class="complete">
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

      <div v-else class="incomplete">
        <div class="name has-text-weight-bold mb-2">{{ torrent.name }}</div>
        <div class="flex status-text">
          <div>{{ prettySize(torrent.size - torrent.completed_bytes) }} left</div>
          <span :title="eta(torrent).toString()">
            ETA:
            <span
              class="primary-text"
              v-if="torrent.down_rate > 0"
            >{{ prettyDate(eta(torrent)) }}</span>
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
  </div>
</template>

<style lang="scss" scoped>
@use "sass:math";
@import "../colors.scss";

.column:nth-child(2) {
  text-align: center;
}
.column:nth-child(3) {
  text-align: right;
}

.app-name {
  cursor: pointer;
}

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

button {
  background-color: $button-color !important;
}

input {
  border-radius: 100px;
}

.status-text {
  color: $text-secondary-color;
  font-size: 14px;
}
.primary-text {
  color: $text-primary-color;
  font-weight: bold;
}

.torrent-card {
  border-radius: 6px;
  background-color: $torrent-unfinished-color;

  &.complete {
    background-color: $torrent-finished-color;
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

.flex {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
