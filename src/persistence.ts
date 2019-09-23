const MSG_DISABLED_LOCALSTORAGE = 'Gitxone uses a function called "localStorage". enable the function of your browser.'

if (!window.localStorage) {
  alert(MSG_DISABLED_LOCALSTORAGE)
}

export const get = (key: string) => {
  return window.localStorage[key]
}

export const set = (key: string, value: any) => {
  window.localStorage.setItem(key, value)
}

export const setData = (data: {[s: string]: any}) => {
  Object.entries(data).map(([k, v]) => set(k, v))
}

export const del = (key: string) => {
  window.localStorage.removeItem(key)
}

export const clear = () => {
  window.localStorage.clear()
}
