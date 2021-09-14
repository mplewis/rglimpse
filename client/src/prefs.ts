type Pref = "accessibleColors";

const Defaults: { [key in Pref]: boolean } = {
  accessibleColors: true,
};

function setDefault(key: Pref): boolean {
  const dfault = Defaults[key];
  set(key, Defaults[key]);
  return dfault;
}

export function set(key: Pref, val: boolean): void {
  localStorage.setItem(key, JSON.stringify(val));
}

export function get(key: Pref): boolean {
  const raw = localStorage.getItem(key);
  if (!raw) {
    return Defaults[key];
  }
  try {
    return JSON.parse(raw);
  } catch (e) {
    return setDefault(key);
  }
}
