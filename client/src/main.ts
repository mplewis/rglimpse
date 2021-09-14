import { createApp } from "vue";
import Pagination from "v-pagination-3";

import App from "./App.vue";
import { store, key } from "./store";

const app = createApp(App);
app.component("pagination", Pagination);
app.use(store, key);
app.mount("#app");
