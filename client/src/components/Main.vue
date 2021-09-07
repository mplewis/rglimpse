<script setup lang="ts">
import { onMounted, ref } from 'vue'
import * as timeago from 'timeago.js';
import en_short from 'timeago.js/lib/lang/en_short'
import bytes from 'bytes'

timeago.register('en_short', en_short)

const props = defineProps<{
  host: string,
  perPage: string,
}>()
const perPage = parseInt(props.perPage)

const loading = ref<boolean>(false)
const error = ref<any>(null)
const page = ref<number>(1)
const total = ref<number>(0)
const torrents = ref<Torrent[]>([])

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

function size(bytez: number): string {
  return bytes(bytez, { unitSeparator: ' ' })
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

async function fetchPage(page: number) {
  // try {
  //   loading.value = true
  //   const count = parseInt(props.perPage)
  //   const offset = (page - 1) * count
  //   const resp = await fetch(`${props.host}/torrents?offset=${offset}&count=${count}`)
  //   const data = await resp.json()
  //   total.value = data.total
  //   torrents.value = data.torrents.map(parse)
  // } catch (e) {
  //   console.log(e)
  //   error.value = e
  // } finally {
  //   loading.value = false
  // }

  torrents.value = [
    {
      completed: false,
      completed_bytes: 10 * 1024 * 1024,
      created: new Date(0),
      down_rate: 100,
      finished: new Date(0),
      hash: 'some-hash-a',
      label: '',
      name: '[SubsPlease] Hamefura S2 - 10 (1080p) [72AC493A].mkv',
      path: '/path/to/my/torrent',
      ratio: 0.5,
      size: 40 * 1024 * 1024,
      started: new Date(Date.now() - 1000 * 60 * 10),
      up_rate: 100,
    },
    {
      completed: true,
      completed_bytes: 40 * 1024 * 1024,
      created: new Date(0),
      down_rate: 100,
      finished: new Date(Date.now() - 1000 * 30),
      hash: 'some-hash-a',
      label: '',
      name: 'La.La.Land.2016.UHD.BluRay.Remux.2160p.HEVC.HDR.Atmos.7.1-HiFi',
      path: '/path/to/my/torrent',
      ratio: 0.5,
      size: 40 * 1024 * 1024,
      started: new Date(Date.now() - 1000 * 60 * 20),
      up_rate: 100,
    },
  ]
  total.value = torrents.value.length
  loading.value = false
}

function diffPage(diff: number) {
  page.value += diff
  fetchPage(page.value)
}

onMounted(() => {
  fetchPage(1)
})

function prettyDate(date: Date): string {
  return timeago.format(date, 'en_short')
}


</script>

<template>
  <div class="mx-3">
    <h1 class="is-size-3 mb-3">rglimpse</h1>
    <div v-if="error">
      <p>Sorry, something went wrong:</p>
      <pre class="my-3"><code>{{ error.toString() }}</code></pre>
      <p>Try refreshing the page.</p>
    </div>

    <div v-else>
      <div class="flex mb-3" v-if="total">
        <button class="button" @click="diffPage(-1)" :disabled="page <= 1">&laquo;</button>
        Page {{ page }} of {{ Math.ceil(total / perPage) }}
        <button
          class="button"
          @click="diffPage(1)"
          :disabled="page >= Math.ceil(total / perPage)"
        >&raquo;</button>
      </div>

      <div v-if="loading">Loading...</div>
      <div
        v-else
        v-for="torrent in torrents"
        :key="torrent.hash"
        :class="{ 'torrent-card': true, 'mb-3': true, 'p-3': true, 'complete': torrent.completed }"
      >
        <div v-if="torrent.completed" class="complete">
          <div class="name has-text-weight-bold">{{ torrent.name }}</div>
          <div class="flex">
            <div class="ratio">Ratio: {{ torrent.ratio }}</div>
            <div>{{ size(torrent.completed_bytes) }} total</div>
            <div :title="torrent.finished.toString()">Done {{ prettyDate(torrent.finished) }}</div>
          </div>
        </div>

        <div v-else class="incomplete">
          <div class="name has-text-weight-bold">{{ torrent.name }}</div>
          <div class="flex">
            <span :title="torrent.started.toString()">{{ prettyDate(torrent.started) }}</span>
            <div>{{ size(torrent.down_rate) }}/s down</div>
            <span :title="eta(torrent).toString()">{{ prettyDate(eta(torrent)) }}</span>
          </div>
          <div class="flex">
            <div>{{ size(torrent.completed_bytes) }} done</div>
            <div>{{ size(torrent.size) }} total</div>
            <div>{{ size(torrent.size - torrent.completed_bytes) }} left</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.torrent-card {
  border-radius: 6px;
  color: white;
  background-color: blue;

  &.complete {
    background-color: green;
  }
}

.name {
  word-break: break-word;
}
</style>


<style lang="scss">
.flex {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
