<script setup lang="ts">
import { ref, watchEffect } from 'vue'
import { useStore } from '../store'
import { Torrent } from '../type'

const store = useStore()

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
    var path = `http://localhost:9081/torrents?offset=${offset}&count=${count}`
    if (query.value.length > 0) {
      path = `http://localhost:9081/torrents?offset=${offset}&count=${count}&query=${query.value}`
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
  <div class="flex">
    <h1 class="app-name is-size-3 mb-3" @click="store.commit('togglePrefs')">rglimpse</h1>
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

    <div v-if="!loading" v-for="torrent in torrents" :key="torrent.hash">
      <TorrentView :torrent="torrent" />
    </div>
  </div>
</template>

<style lang="scss">
@import "../colors.scss";

.app-name {
  cursor: pointer;
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

.flex {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
