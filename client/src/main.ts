import { createApp } from "vue";
import Pagination from "v-pagination-3";

import { store, key } from "./store";
import App from "./App.vue";
import TorrentView from "./components/TorrentView.vue";

const app = createApp(App);
app.component("pagination", Pagination);
app.component("TorrentView", TorrentView);
app.use(store, key);
app.mount("#app");
