// 永続化データを取り扱うための関数群

const MSG_DISABLED_LOCALSTORAGE = 'gitxone uses a function called "localStorage". enable the function of your browser.'

if (!window.localStorage) {
  alert(MSG_DISABLED_LOCALSTORAGE)
}

export const get = (key, json=true) => {
  const value = window.localStorage[key]
  if (!json) {
    return value
  }
  try {
    return JSON.parse(value)
  } catch {
    return null
  }
}

export const set = (key, value, json=true) => {
  window.localStorage.setItem(key, json ? JSON.stringify(value) : value)
}

export const setData = (data) => {
  Object.entries(data).map(([k, v]) => set(k, v))
}

export const del = (key) => {
  window.localStorage.removeItem(key)
}

export const clear = () => {
  window.localStorage.clear()
}
