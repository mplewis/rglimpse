import { InjectionKey } from "vue";
import { Store, createStore, useStore as _useStore } from "vuex";
import { set, get } from "./prefs";

export const key: InjectionKey<Store<State>> = Symbol();

export interface State {
  prefsOpen: boolean;
  accessibleColors: boolean;
}

export const store = createStore<State>({
  state() {
    return {
      prefsOpen: false,
      accessibleColors: get("accessibleColors"),
    };
  },
  mutations: {
    togglePrefs(state) {
      state.prefsOpen = !state.prefsOpen;
    },
    setAccessibleColors(state, value: boolean) {
      state.accessibleColors = value;
      set("accessibleColors", value);
    },
  },
});

export function useStore(): Store<State> {
  return _useStore(key);
}
