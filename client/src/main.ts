import { createApp } from "vue";
import App from "./App.vue";
import Pagination from "v-pagination-3";

const app = createApp(App);
app.component("pagination", Pagination);
app.mount("#app");
