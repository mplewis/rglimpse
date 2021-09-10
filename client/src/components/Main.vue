<script setup lang="ts">
import { onMounted, ref, watchEffect } from 'vue'
import { prettyDate, prettySize } from '../logic/pretty'

const props = defineProps<{
  host: string,
  perPage: string,
}>()
const perPage = parseInt(props.perPage)

const loading = ref<boolean>(true)
const error = ref<any>(null)
const page = ref<number>(1)
const total = ref<number>(0)
const torrents = ref<Torrent[]>([])
const query = ref<string>()
const queryRaw = ref<string>()

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

const wip = {
  completed: false,
  completed_bytes: 291.56 * 1024 * 1024,
  created: new Date(0),
  down_rate: 75400,
  finished: new Date(0),
  hash: 'some-hash-a',
  label: '',
  name: '[SubsPlease] Hamefura S2 - 10 (1080p) [72AC493A].mkv',
  path: '/path/to/my/torrent',
  ratio: 0.108,
  size: 378.812 * 1024 * 1024,
  started: new Date(Date.now() - 1000 * 60 * 10),
  up_rate: 100,
}
const done = {
  completed: true,
  completed_bytes: 18.32 * 1024 * 1024 * 1024,
  created: new Date(0),
  down_rate: 100,
  finished: new Date(Date.now() - 1000 * 30),
  hash: 'some-hash-a',
  label: '',
  name: 'La.La.Land.2016.UHD.BluRay.Remux.2160p.HEVC.HDR.Atmos.7.1-HiFi',
  path: '/path/to/my/torrent',
  ratio: 0.081,
  size: 40 * 1024 * 1024,
  started: new Date(Date.now() - 1000 * 60 * 20),
  up_rate: 128.52 * 1024,
}

// TODO: Auto refresh
async function fetchPage(page: number) {
  try {
    loading.value = true
    const count = parseInt(props.perPage)
    const offset = (page - 1) * count
    const resp = await fetch(`${props.host}/torrents?offset=${offset}&count=${count}`)
    const data = await resp.json()
    total.value = data.total
    torrents.value = data.torrents.map(parse)
  } catch (e) {
    console.log(e)
    error.value = e
  } finally {
    loading.value = false
  }

  // torrents.value = [wip, wip, wip]
  // for (let x = 0; x < 7; x++) {
  //   torrents.value.push(done)
  // }
  // total.value = torrents.value.length
  // loading.value = false
}

function diffPage(diff: number) {
  page.value += diff
  fetchPage(page.value)
}

onMounted(() => {
  fetchPage(1)
})

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
    <h1 class="appname is-size-3 mb-3" @click="showSettings = true">rglimpse</h1>
    <p v-if="loading">Loading...</p>
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
      {{ (page - 1) * perPage + 1 }}–{{ page * perPage }} of {{ total }}
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

.appname {
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
