import { InjectionKey } from "vue";
import { Store, createStore, useStore as _useStore } from "vuex";

export const key: InjectionKey<Store<State>> = Symbol();

export interface State {
  prefsOpen: boolean;
}

export const store = createStore<State>({
  state() {
    return {
      prefsOpen: false,
    };
  },
  mutations: {
    togglePrefs(state) {
      state.prefsOpen = !state.prefsOpen;
    },
  },
});

export function useStore(): Store<State> {
  return _useStore(key);
}
